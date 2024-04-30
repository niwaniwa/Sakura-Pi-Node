<div align="center">
<img src="https://raw.githubusercontent.com/niwaniwa/Project-Sakura/main/Assets/icon.png" alt="Sakura logo" width="400"/>

** 🔑🌸🔑 [Sakura](https://github.com/niwaniwa/Project-Sakura) is a simple access control system. **

</div>

# 🌸 Sakura-Pi-Node

Hardware node for RaspberryPi

Project Page **[Sakura](https://github.com/niwaniwa/Project-Sakura)**

## 💉 Dependence
- Raspberry Pi
- [Pasori](https://github.com/bamchoh/pasori)
- [go-rpio](https://github.com/stianeikeland/go-rpio)
- libusb

## How to use

### build from source

1. Installation of Dependencies
```bash
apt install -y gcc-arm-linux-gnueabihf libusb-1.0-0-dev libusb-1.0-0-dev:armhf
```

2. Build from source
```bash
go mod tidy
go build -o Sakura-Pi-Node ./cmd
```

3. Run
Administrative privileges are required. (for Pasori)
```bash
sudo ./main
```

### build by docker compose
Automatically run in privileged mode.
```bash
docker compose up -d
```

## License
MIT License