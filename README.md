# Evade
Divide the file into blocks according to bytes.

中文：将文件按照字节进行分块．

##Instructions for use(使用说明)
```
Usage:
  main <-i filename> [-o output] [-p prefix] [-s 4096] [--disable-append]

Application Options:
  -i, --input-file=File           Input source file
  -o, --output-dir=Dir            Output directory (default: output)
  -p, --filename-prefix=Prefix    Output file name prefix (default: InputFileName)
  -s, --size=Byte                 Shard size (default: 4096)
      --disable-append            Disable append mode

Help Options:
  -h, --help                      Show this help message


```
