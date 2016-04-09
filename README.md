# AGI Server


```
Set(STARTVAR=VAR) variable in your dialplan before exec AGI(agi://127.0.0.1/AGIServer)

available VARS:
    block - for blocking ip addresses via iptables (testing)
    inbound - for incoming calls
    CONFBRIDGES:
        confbridge_channelredirect - for confbridge create as admin and add callee (bridged channel)
        confbridge_access - confbridge access
        confbridge_addmembers - try to add member
        confbridge_confs - add member
```

## Example Internal Confbridge:
```
[macro-inv]
        exten => s,1,Set(STARTVAR=confbridge_channelredirect)
        exten => s,n,AGI(agi://0.0.0.0/AGIServer)

[dyn-nway]
        exten => _X.,1,Set(STARTVAR=confbridge_access)
        exten => _X.,n,AGI(agi://0.0.0.0/AGIServer)

[conf_add]
        exten => s,1,Set(STARTVAR=confbridge_addmembers)
        exten => s,n,AGI(agi://0.0.0.0/AGIServer)

[confs]
        exten => _X.,1,Set(STARTVAR=confbridge_confs)
        exten => _X.,n,AGI(agi://0.0.0.0/AGIServer)
```# gotcpd
