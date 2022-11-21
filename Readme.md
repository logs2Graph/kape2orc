# KAPE2ORC

DFIR-ORC is an Open-Source Artefact Collection tool develop by ANSSI (French Cybersecurity Nationnal Agency). This tool can be configured to collect artefact like Kape and to embed use third party tools to execute them on a live machine.

Howerever, the community around this tool is smaller than Kape's and therefore the collection of artefact collected by the default configuration is not as large and as supported as Kape's. This project aims to provide a way to convert Kape's configuration file in DFIR-ORC configuration file. This way, the work done on Kape can be reused in DFIR-ORC. 

## Usage

```
Usage of ./kape2orc:
Usage of /tmp/go-build2394873438/b001/exe/master:
  -kape string
    	Directory where Tkape files are located (default "./kape")
  -master string
    	Master Tkape file (default "./kape/Compound/!SANS_Triage.tkape")
  -orc string
    	Output Directory to write Orc config files (default "./orc")
  -verbose
    	Verbose mode
```

DFIR-ORC doesn't support recursive dependencies. So we need to assign a master file and the dependencies will be flatten based on the master Kape file provided.

## Note

For your information, DFIR ORC doesn't handle recursive dependencies. So when the tool encounters a recursive dependency, it will break down the dependency and merge it in a single File.

## Todo

- [ ] Handle several Master files
