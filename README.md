# Tmux workspace

Tạo tab tmux, rename với json config.

```json
{
  "default": [
    {
      "path": "/home/ubuntu/code/learning/dsa/algorithm",
      "name": "algo"
    },
    {
      "path": "/home/ubuntu/code/learning/backend/golang",
      "name": "go"
    }
  ],
  "backend": [
    {
      "path": "/home/ubuntu/code/learning/backend/golang",
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

- Go user

```sh
go install github.com/tooltoys/tmux-session@latest
```

3. Usage

```sh
tmux-session -name=default
```
