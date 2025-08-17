%%{init: {'theme': 'base', 'themeVariables': { 'primaryColor': '#ffd8d8', 'primaryBorderColor': '#ff9e9e', 'lineColor': '#ff9e9e'}}}%%
graph TD
    subgraph "Wish Server"
        SSH[SSH Server] -->|Hosts| TUI
        TUI[BubbleTea TUI] -->|Uses| Redis[(Redis)]
    end
    
    User1[User 1] -->|SSH Connection| SSH
    User2[User 2] -->|SSH Connection| SSH
    
    TUI -->|Publishes| Redis
    TUI -->|Subscribes| Redis