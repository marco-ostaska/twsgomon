<img src="logo.png" alt="twsgomon" title="twsgomon" align="right" width="160">

# twsgom⌕n - TWS monitoring 

[![Go Report Card](https://goreportcard.com/badge/github.com/marco-ostaska/twsgomon)](https://goreportcard.com/report/github.com/marco-ostaska/twsgomon)

`twsgom⌕n` is a tool that monitors TWS event.log to check and parse for job abends, job fails, job late, etc with a feature to send webhook to Netcool MessageBus

## Installation

Just download and run the binary twsgomon file from whenever you want using `-config` and provide and json configuration file

## Configurations files 

### twsgomon.json 

- `evenlog.path` path to event.log to monitor
- `readFromBeginning` whether it should read all event.log or from the moment twsgomon started
- `interval` interval cycle it should checks fir updated in event.log
- `twsgomonlog.path` where it should save twsgomon log
- `debug.level` being debug the higher debug level information and fatal lowest debug level information. 
  - The debug levels are: (`debug`, `info`, `warn`, `error`, `fatal`) 
- `alert.config`: json configuration path to configure what should alert

#### example: 

```json
    "evenlog.path" : "/home/marcoan/go/src/github.com/marco-ostaska/twsgomon/event.log",
    "readFromBeginning" : true,
    "interval": 2,
    "twsgomonlog.path" : "/tmp/twsgomon.log",
    "debug.level" : "info",
    "alert.config": "/home/marcoan/go/src/github.com/marco-ostaska/twsgomon/internal/configs/alerts.json"
}
```


### alerts.json 

This file is mounted using [*TWS Appendix A. Job scheduling events format documentation*](http://publib.boulder.ibm.com/tividd/td/TWS/SC32-1276-02/en_US/HTML/plusmst70.htm) as base




- `alerts` an array of alerts 
- `eventNumber` is the Event Number mentioned on  [*TWS Appendix A. Job scheduling events format documentation*](http://publib.boulder.ibm.com/tividd/td/TWS/SC32-1276-02/en_US/HTML/plusmst70.htm)
- `PositionalFields` is an array, for mapping see  [*TWS Appendix A. Job scheduling events format documentation*](http://publib.boulder.ibm.com/tividd/td/TWS/SC32-1276-02/en_US/HTML/plusmst70.htm)
- `MessageBus` bool param that tells twsgomin if it should send webhook to NetcoolMessageBus
- `AlertKey` the alertKey Netcool should receive
- `Severity` Netcool severity 
- `MessageBusURL` NetcoolMessageBus URL

#### example: 

```json
{
  "alerts": [
      {
        "eventNumber":101,
        "PositionalFields": ["1:um", "4:quatro"],
        "MessageBus" : true,
        "AlertKey" : "twsgomon",
        "Severity": "2",
        "MessageBusURL": "http://localhost:31311"
       },
       {
        "eventNumber":103,
        "PositionalFields": ["1:um", "4:quatro"]
       }
    ]
}

```
