package main

import (
	"encoding/xml"
	"flag"
	//"fmt"
	yaml "gopkg.in/yaml.v3"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

var (
	source      = flag.String("kape", "./kape", "Directory where Tkape files are located")
	output      = flag.String("orc", "./orc", "Output Directory to write Orc config files")
	master      = flag.String("master", "./kape/Compound/!SANS_Triage.tkape", "Master Tkape file")
	keep_unused = flag.Bool("keep_unused", false, "Keep unused files in the output directory")
	verbose     = flag.Bool("verbose", false, "Verbose mode")
)

func handleErr(err error) {
	if err != nil {
		log.Println(err)
	}
}

type Target struct {
	Name      string `yaml:"Name"`
	Category  string `yaml:"Category"`
	Path      string `yaml:"Path"`
	FileMask  string `yaml:"FileMask"`
	Recursive bool   `yaml:"Recursive"`
}

type KapeFile struct {
	Path        string
	Name        string
	Description string   `yaml:"Description"`
	Author      string   `yaml:"Author"`
	Version     string   `yaml:"Version"`
	Targets     []Target `yaml:"Targets"`
}

type OrcConfig struct {
	XMLName   xml.Name `xml:"getthis"`
	Text      string   `xml:",chardata"`
	Reportall string   `xml:"reportall,attr,omitempty"`
	Hash      string   `xml:"hash,attr,omitempty"`
	Output    struct {
		Text        string `xml:",chardata"`
		Compression string `xml:"compression,attr,omitempty"`
	} `xml:"output"`
	Location string  `xml:"location"`
	Samples  Samples `xml:"samples"`
}

type Samples struct {
	Text              string   `xml:",chardata"`
	MaxPerSampleBytes string   `xml:"MaxPerSampleBytes,attr,omitempty"`
	MaxTotalBytes     string   `xml:"MaxTotalBytes,attr,omitempty"`
	MaxSampleCount    string   `xml:"MaxSampleCount,attr,omitempty"`
	Sample            []Sample `xml:"sample"`
}

type Sample struct {
	Text              string     `xml:",chardata"`
	Name              string     `xml:"name,attr,omitempty"`
	MaxPerSampleBytes string     `xml:"MaxPerSampleBytes,attr,omitempty"`
	NtfsFind          []NTFSFind `xml:"ntfs_find"`
}

type NTFSFind struct {
	Text      string `xml:",chardata"`
	PathMatch string `xml:"path_match,attr,omitempty"`
	Name      string `xml:"name,attr,omitempty"`
}

type Log struct {
	Text        string `xml:",chardata"`
	Disposition string `xml:"disposition,attr,omitempty"`
}

type Outline struct {
	XMLName     xml.Name `xml:"outline"`
	Text        string   `xml:",chardata"`
	Disposition string   `xml:"disposition,attr,omitempty"`
}

type Restriction struct {
	Text             string `xml:",chardata"`
	ElapsedTimeLimit string `xml:"ElapsedTimeLimit,attr,omitempty"`
}

type Output struct {
	Text     string `xml:",chardata"`
	Name     string `xml:"name,attr,omitempty"`
	Source   string `xml:"source,attr,omitempty"`
	Argument string `xml:"argument,attr,omitempty"`
}

type Command struct {
	Text       string `xml:",chardata"`
	Keyword    string `xml:"keyword,attr,omitempty"`
	Optional   string `xml:"optional,attr,omitempty"`
	Queue      string `xml:"queue,attr,omitempty"`
	Winver     string `xml:"winver,attr,omitempty"`
	Systemtype string `xml:"systemtype,attr,omitempty"`
	Execute    struct {
		Text string `xml:",chardata"`
		Name string `xml:"name,attr"`
		Run  string `xml:"run,attr"`
	} `xml:"execute"`
	Argument []string `xml:"argument"`
	Output   []Output
}

type Archive struct {
	Text           string        `xml:",chardata"`
	Name           string        `xml:"name,attr,omitempty"`
	Keyword        string        `xml:"keyword,attr,omitempty"`
	Format         string        `xml:"format,attr,omitempty"`
	File           []File        `xml:"file,omitempty"`
	Concurrency    string        `xml:"concurrency,attr,omitempty"`
	Repeat         string        `xml:"repeat,attr,omitempty"`
	Compression    string        `xml:"compression,attr,omitempty"`
	ArchiveTimeout string        `xml:"archive_timeout,attr,omitempty"`
	Optional       string        `xml:"optional,attr,omitempty"`
	Restrictions   []Restriction `xml:"restrictions,omitempty"`
	Command        []Command     `xml:"command"`
}

type Wolf struct {
	XMLName        xml.Name  `xml:"wolf"`
	Text           string    `xml:",chardata"`
	Childdebug     string    `xml:"childdebug,attr,omitempty"`
	CommandTimeout string    `xml:"command_timeout,attr,omitempty"`
	Log            Log       `xml:"log"`
	Outline        Outline   `xml:"outline"`
	Archive        []Archive `xml:"archive"`
}

type File struct {
	Text string `xml:",chardata"`
	Name string `xml:"name,attr"`
	Path string `xml:"path,attr"`
}

type Toolembed struct {
	XMLName xml.Name `xml:"toolembed"`
	Text    string   `xml:",chardata"`
	Input   string   `xml:"input"`
	Output  string   `xml:"output"`
	Run64   struct {
		Text string `xml:",chardata"`
		Args string `xml:"args,attr"`
	} `xml:"run64"`
	Run32 struct {
		Text string `xml:",chardata"`
		Args string `xml:"args,attr"`
	} `xml:"run32"`
	Archive Archive `xml:"archive"`
	File    []File  `xml:"file"`
}

func ParseKape(filename string) KapeFile {
	var kapefile KapeFile
	data, err := ioutil.ReadFile(filename)
	handleErr(err)
	err = yaml.Unmarshal(data, &kapefile)
	handleErr(err)

	//Get filename
	splitted_path := strings.Split(filename, "/")
	handleErr(err)
	kapefile.Name = splitted_path[len(splitted_path)-1]
	kapefile.Name = strings.Replace(kapefile.Name, ".tkape", "", -1)

	return kapefile
}

func ConvertGetThis(kapefile KapeFile) []byte {
	var orcConfig OrcConfig

	orcConfig.Location = "%SystemDrive%"

	for _, target := range kapefile.Targets {
		var sample Sample
		var ntfsFind NTFSFind
		sample.Name = strings.Replace(target.Name, " ", "_", -1)
		path := strings.Replace(target.Path, "C:\\", "\\", 1)
		path = strings.Replace(path, "%user%", "*", -1)

		// Merge FileMasks and Path
		if target.FileMask != "" && target.Path != "" {
			ntfsFind.PathMatch = path + target.FileMask
		} else if target.FileMask != "" {
			ntfsFind.Name = target.FileMask
		} else {
			ntfsFind.PathMatch = path + "*"
		}

		sample.NtfsFind = append(sample.NtfsFind, ntfsFind)
		orcConfig.Samples.Sample = append(orcConfig.Samples.Sample, sample)
	}

	str, err := xml.MarshalIndent(orcConfig, "", "	")
	handleErr(err)

	return str
}

func ConvertWolf(kapefile KapeFile) []byte {
	var wolfConfig Wolf
	var archive Archive
	archive.Name = kapefile.Name + ".7z"
	archive.Keyword = kapefile.Name
	archive.Compression = "fast"
	archive.Repeat = "Once"
	archive.Concurrency = "4"
	wolfConfig.Childdebug = "true"
	wolfConfig.CommandTimeout = "3600"
	wolfConfig.Log.Text = "DFIR-ORC_{SystemType}_{FullComputerName}_{TimeStamp}.log"
	wolfConfig.Outline.Text = "DFIR-ORC_{SystemType}_{FullComputerName}_{TimeStamp}.json"
	wolfConfig.Outline.Disposition = "truncate"

	for _, target := range kapefile.Targets {
		var command Command
		var file_out Output
		var file_err Output

		command.Keyword = target.Name
		command.Execute.Name = "Orc.exe"
		command.Execute.Run = "self:#GetThis"
		target_path := target.Path
		target_path = strings.Replace(target_path, ".tkape", "_config.xml", 1)
		command.Argument = append(command.Argument, "/config=res:#"+target_path+" /NoLimits")
		file_out.Name = target.Name + ".7z"
		file_out.Source = "File"
		file_out.Argument = "/out={FileName}"
		file_err.Source = "StdOutErr"
		file_err.Name = target.Name + ".log"

		command.Output = append(command.Output, file_out)
		command.Output = append(command.Output, file_err)
		archive.Command = append(archive.Command, command)

	}

	wolfConfig.Archive = append(wolfConfig.Archive, archive)

	orcstring, err := xml.MarshalIndent(wolfConfig, "", "	")
	handleErr(err)
	return orcstring
}

// Get Output Path from Source File Path and Name
func GetOutputPath(path string) string {
	splitted_path := strings.Split(path, "/")
	splitted_source := strings.Split(*source, "/")
	j := 0
	// We compare the source path and the target path to find the common part
	for i := 0; i < len(splitted_path) && i < len(splitted_source); i++ {
		j = i
		continue
	}
	// We remove the common part from the target path and replace it by the output path
	var splitted_output []string
	splitted_output = append(splitted_output, *output)
	splitted_output = append(splitted_output, splitted_path[j+1:]...)

	// We replace the extension by .xml
	splitted_output[len(splitted_output)-1] = strings.Replace(splitted_output[len(splitted_output)-1], ".tkape", "_config.xml", 1)

	return strings.Join(splitted_output, "/")
}

func Export(kapefiles []KapeFile) {
	for _, kapefile := range kapefiles {
		var data []byte
		output_path := GetOutputPath(kapefile.Path)
		if IsGetThis(kapefile) {
			data = ConvertGetThis(kapefile)
		} else if IsCompound(kapefile) {
			data = ConvertWolf(kapefile)
		} else {
			if *verbose {
				log.Println("Error: Failed to Export " + kapefile.Name)
				log.Println("Error: ", kapefile.Targets)
			}
		}

		if len(data) > 0 {
			// We parse the output path the create the directory if it doesn't exist
			splitted_path := strings.Split(output_path, "/")
			err := os.MkdirAll(strings.Join(splitted_path[:len(splitted_path)-1], "/"), 0755)
			handleErr(err)

			// We write the file
			err = os.WriteFile(output_path, data, 0644)
			handleErr(err)
		}
	}
}

func IsCompound(kapefile KapeFile) bool {
	hasTkape := false
	hasGetThis := false
	for _, target := range kapefile.Targets {
		if strings.Contains(target.Path, ".tkape") {
			hasTkape = true
		} else {
			hasGetThis = true
		}
	}
	return hasTkape && !hasGetThis
}

func IsMixed(kapefile KapeFile) bool {
	hasTkape := false
	hasGetThis := false
	for _, target := range kapefile.Targets {
		if strings.Contains(target.Path, ".tkape") {
			hasTkape = true
		} else {
			hasGetThis = true
		}
	}
	return hasTkape && hasGetThis
}

func IsGetThis(kapefile KapeFile) bool {
	hasTkape := false
	hasGetThis := false
	for _, target := range kapefile.Targets {
		if strings.Contains(target.Path, ".tkape") {
			hasTkape = true
		} else {
			hasGetThis = true
		}
	}
	return !hasTkape && hasGetThis
}

func GenerateEmbed(files []KapeFile, master KapeFile) Toolembed {
	var toolEmbed Toolembed
	toolEmbed.Input = ".\\tools\\DFIR-Orc_x86.exe"
	toolEmbed.Output = ".\\output\\%ORC_OUTPUT%"
	toolEmbed.Run64.Args = "WolfLauncher"
	toolEmbed.Run64.Text = "7z:#Tools|DFIR-Orc_x64.exe"
	toolEmbed.Run32.Args = "WolfLauncher"
	toolEmbed.Run32.Text = "self:#"

	var archive Archive
	archive.Name = "Tools"
	archive.Format = "7z"
	archive.Compression = "Ultra"

	var file File
	file.Name = "DFIR-Orc_x64.exe"
	file.Path = ".\\tools\\DFIR-Orc_x64.exe"

	archive.File = append(archive.File, file)

	toolEmbed.Archive = archive

	var master_file File
	master_file.Name = "WOLFLAUNCHER_CONFIG"
	splitted_path := strings.Split(master.Path, "/")
	if splitted_path[0] == "." {
		splitted_path[1] = "%ORC_CONFIG_FOLDER%"
	} else {
		splitted_path[0] = "%ORC_CONFIG_FOLDER%"
	}
	splitted_path[len(splitted_path)-1] = strings.Replace(splitted_path[len(splitted_path)-1], ".tkape", "_config.xml", 1)
	master_file.Path = strings.Join(splitted_path, "/")

	toolEmbed.File = append(toolEmbed.File, master_file)

	for _, f := range files {
		file := File{}
		splitted_path := strings.Split(f.Path, "/")

		if splitted_path[0] == "." {
			splitted_path[1] = "%ORC_CONFIG_FOLDER%"
		} else {
			splitted_path[0] = "%ORC_CONFIG_FOLDER%"
		}

		splitted_path[len(splitted_path)-1] = strings.Replace(splitted_path[len(splitted_path)-1], ".tkape", "_config.xml", 1)

		file.Path = strings.Join(splitted_path, "\\")
		file.Name = splitted_path[len(splitted_path)-1]

		toolEmbed.File = append(toolEmbed.File, file)
	}

	return toolEmbed
}

// Parse all kape file in the source directory except the master kape file
func ParseKapeDirectory(path string, master string) []KapeFile {
	files, err := os.ReadDir(path)
	handleErr(err)
	kapefiles := []KapeFile{}
	for _, f := range files {
		file_path := path + "/" + f.Name()
		if file_path == master {
			continue
		}
		if strings.Contains(f.Name(), ".tkape") {
			kapefile := ParseKape(file_path)
			kapefile.Path = path + "/" + f.Name()
			kapefile.Name = strings.Replace(f.Name(), ".tkape", "", 1)
			kapefiles = append(kapefiles, kapefile)
		} else if f.Type().IsDir() {
			//log.Println("Directory found")
			tmp_kape := ParseKapeDirectory(file_path, master)
			kapefiles = append(kapefiles, tmp_kape...)
		}
	}
	return kapefiles
}

// Split Mixed Kapefile into GetThis and Compound
// Also works on GetThis compatible files.
func SplitKape(kapefile KapeFile) (KapeFile, KapeFile) {
	var compound KapeFile
	var getthis KapeFile

	compound.Name = kapefile.Name
	compound.Path = kapefile.Path
	compound.Description = kapefile.Description
	compound.Author = kapefile.Author
	compound.Version = kapefile.Version

	getthis.Name = kapefile.Name + "_getthis"
	getthis.Path = strings.Replace(kapefile.Path, ".tkape", "_getthis.tkape", 1)
	getthis.Description = kapefile.Description
	getthis.Author = kapefile.Author
	getthis.Version = kapefile.Version

	for _, target := range kapefile.Targets {
		if strings.Contains(target.Path, ".tkape") {
			compound.Targets = append(compound.Targets, target)
		} else {
			getthis.Targets = append(getthis.Targets, target)
		}
	}

	// Add Target to Compound to reference GetThis
	var getthis_target Target
	getthis_target.Name = getthis.Name
	getthis_target.Path = getthis.Name + ".tkape"
	compound.Targets = append(compound.Targets, getthis_target)

	return compound, getthis
}

func Flatten(kapefiles []KapeFile, toFlatten KapeFile) KapeFile {
	toDelete := []int{}
	for i, target := range toFlatten.Targets {
		if strings.Contains(target.Path, ".tkape") {
			found := false
			for _, kapefile := range kapefiles {
				if kapefile.Name+".tkape" == target.Path {
					if IsMixed(kapefile) {
						log.Println("Error: Flatten Mixed Kapefile : ", kapefile.Name)
					}

					found = true
					var flatten KapeFile
					flatten = Flatten(kapefiles, kapefile)
					toFlatten.Targets = append(toFlatten.Targets, flatten.Targets...)
					toDelete = append(toDelete, i)
					/*
						if target.Path == "ApplicationEvents.tkape" {
							log.Println("Info: Found ApplicationEvents for " + toFlatten.Path + " at " + kapefile.Path)
							log.Println("Info: Targets : ", toFlatten.Targets)
						}
					*/
				}
			}
			if !found {
				if *verbose {
					log.Println("Warn: could not find " + target.Path + " for flattening " + toFlatten.Name + ". Skipping...")
				}
				toDelete = append(toDelete, i)
			}
		}
	}

	// Remove flattened targets
	for i := len(toDelete) - 1; i >= 0; i-- {
		toFlatten.Targets = append(toFlatten.Targets[:toDelete[i]], toFlatten.Targets[toDelete[i]+1:]...)
	}

	return toFlatten
}

// We replace all Compounds by GetThis compatible files
func ConvertCompound(kapefiles []KapeFile) []KapeFile {
	for i, kapefile := range kapefiles {
		if !IsGetThis(kapefile) {
			kapefiles[i] = Flatten(kapefiles, kapefile)
		}
	}

	return kapefiles
}

func GetUsedKapefile(kapefiles []KapeFile, kapefile KapeFile) []KapeFile {
	used_kapefiles := []KapeFile{}
	for _, target := range kapefile.Targets {
		if strings.Contains(target.Path, ".tkape") {
			for _, k := range kapefiles {
				if k.Name+".tkape" == target.Path {
					used_kapefiles = append(used_kapefiles, k)
				}
			}
		}
	}
	return used_kapefiles
}

func main() {
	var toolEmbed Toolembed

	flag.Parse()

	master_kape := ParseKape(*master)
	master_kape.Path = *master

	if IsCompound(master_kape) {
		kapefiles := ParseKapeDirectory(*source, *master)
		kapefiles = ConvertCompound(kapefiles)

		if *keep_unused {
			Export(kapefiles)
			toolEmbed = GenerateEmbed(kapefiles, master_kape)

		} else {
			used_kapefiles := GetUsedKapefile(kapefiles, master_kape)
			Export(used_kapefiles)
			toolEmbed = GenerateEmbed(used_kapefiles, master_kape)
		}

		Export([]KapeFile{master_kape})

		// Generate the embed file
		data, err := xml.MarshalIndent(toolEmbed, "", "  ")
		handleErr(err)
		err = ioutil.WriteFile(*output+"/DFIR-ORC_embed.xml", data, 0644)
		handleErr(err)
	} else {
		master_compound, master_getthis := SplitKape(master_kape)
		kapefiles := ParseKapeDirectory(*source, *master)
		kapefiles = append(kapefiles, master_getthis)
		kapefiles = ConvertCompound(kapefiles)

		if *keep_unused {
			Export(kapefiles)
			toolEmbed = GenerateEmbed(kapefiles, master_compound)
		} else {
			used_kapefiles := GetUsedKapefile(kapefiles, master_compound)
			Export(used_kapefiles)
			Export([]KapeFile{master_compound})
			toolEmbed = GenerateEmbed(used_kapefiles, master_compound)
		}

		// Generate the embed file
		data, err := xml.MarshalIndent(toolEmbed, "", "  ")
		handleErr(err)
		err = ioutil.WriteFile(*output+"/DFIR-ORC_embed.xml", data, 0644)
		handleErr(err)
	}

}
