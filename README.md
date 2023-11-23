# Tmux workspace

Tạo tab tmux, rename với json config.

```json
{
  "default": [
    {
      "path": "/home/ubuntu/algorithm",
      "name": "algo"
    },
    {
      "path": "/home/ubuntu/music",
      "name": "go"
    }
  ],
  "backend": [
    {
      "path": "/home/ubuntu/backend",
      "name": "go"
    }
  ]
}
```

Trong đó 

+ default: tên session
+ path: là đường dẫn đến dir cần thao tác.
+ name: custom name cho tmux panel. 

## Installation

1. Tạo file ~/.sessionrc 

2. Download 


```sh
# Gopher
go install github.com/tooltoys/tmux-session@v1.0.1
```

3. Usage

```sh
tmux-session -name=default
```
