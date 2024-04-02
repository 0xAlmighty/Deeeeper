# Deeeeper 🕵️‍♂️📱

Welcome to **Deeeeper**, the ultimate tool for going deep into the heart of Android applications. Whether you're a **penetration tester**, a **security enthusiast**, or a **curious developer**, Deeeeper offers you the capability to decompile APK files and uncover hidden activities and deeplinks. 
## 🌟 Features

- **Decompile APKs:** Convert APK files back into understandable source code.
- **Extract Components:** Pull out activities, services, receivers, and their intents with ease.
- **Deeplink Discovery:** Identify and construct deeplink URIs to understand how apps communicate internally and with the outside world.
- **Colorful Console Output:** Because who doesn't like a bit of color in their terminal?

**Requirements**

- Go 1.22.0 or later
- apktool

Get started with Deeeeper in just a few steps:

Download the release from the release page

1. Clone the repository
```
git clone https://github.com/0xAlmighty/deeeeper.git
```
2. Change to the git directory
```
cd deeeeper
```
3. Build
```
go build
```
## 🚀 Usage

Run Deeeeper with the following commands:

To **decompile an APK** and analyze its components:

```
./deeeeper -apk path/to/your/app.apk
```

If APK is already decompiled, target the folder:

```
./deeeeper -folder path/to/your/folder
```

Need **help**? Just ask:

```shell
./deeeeper --help

Usage: deeeeper [OPTIONS]
Options:
  -apk <path>       Path to the APK file to be decompiled
  -folder <path>    Folder to search in if APK is already decompiled
  -h, --help        Display this help and exit
```
## 🤝 Contributing

Your contributions are what make the magic happen! Whether you're fixing bugs, adding new features, or improving documentation, we welcome your pull requests and issues. Feel free to reach out.
## 📜 License

Deeeeper is made available under the MIT License. For more details, see the LICENSE file.
## 📢 Let's Make Something Great Together!

I'm excited to see how Deeeeper evolves with your input. Don't hesitate to reach out, propose changes, or simply star 🌟 the project if you find it useful. 

Happy hacking!
