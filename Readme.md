# KAPE2ORC

## Usage

```
Usage of ./kape2orc:
  -kape string
    	Directory where Tkape files are located (default "./kape")
  -orc string
    	Output Directory to write Orc config files (default "./orc")
  -verbose
    	Verbose mode
```

## Note

For your information, DFIR ORC doesn't handle recursive dependencies. So when the tool encounters a recursive dependency, it will break down the dependency and merge it in a single File.


