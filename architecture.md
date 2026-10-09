```mermaid
%%{init: {'theme': 'base', 'themeVariables': { 'primaryColor': '#d8ffd8', 'primaryBorderColor': '#9eff9e', 'lineColor': '#9eff9e'}}}%%
graph TD
    subgraph "Terminal Chat"
        TUI[BubbleTea TUI] -->|Subscribes/Publishes| Broker[In-Memory Broker]
        TUI -->|Optional TCP| TCP[TCP Server/Client]
        TCP -->|Relays JSON| Broker
        TUI -->|Disconnect| TCP
    end

    User1[Local User] -->|Terminal| TUI
    User2[Remote User] -->|TCP Connection| TCP
    UserN[Remote User N] -->|TCP Connection| TCP

    TCP -->|Broadcasts| RemoteClients[Connected Remote Clients]
```