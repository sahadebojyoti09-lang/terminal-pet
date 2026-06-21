# 🐾 Puchi - The Ultimate Terminal Pet

A blazing fast, ultra-smooth virtual terminal pet built with **Go** and the **Bubble Tea** TUI framework. Puchi lives inside your terminal, blinks, gets hungry, dances, and dispenses random wisdom using `fortune`.

---

## ✨ Features

* **Zero Flicker:** Uses double-buffering terminal rendering for buttery-smooth animations.
* **Low Footprint:** Compiles into a single native binary with near 0% CPU idle usage.
* **Fortune-Powered:** Pressing `s` makes Puchi speak short quotes directly from your system's `fortune` database.
* **Persistent Speech:** Text wraps safely downward and freezes on screen so you can read fortunes completely at your own pace without jittering the layout walls.

---

## 🐧 Linux Installation

### 1. Prerequisites

Make sure you have `go` and the `fortune` database installed on your system:

* **Arch Linux:**
```
  sudo pacman -S go fortune-mod
```
* **Debian / Ubuntu:**

```
  sudo apt update
```
```
  sudo apt install golang-go fortune-mod fortunes-min -y
```

2. Clone & Run
Clone the repository:

```
git clone [https://github.com/sahadebojyoti09-lang/terminal-pet.git](https://github.com/sahadebojyoti09-lang/terminal-pet.git)
```
Move into the project directory:

```
cd terminal-pet
```
Build the optimized binary:

```
go build -o pet main.go
```
Move it to your global binaries path so it can be run from anywhere (this will request your sudo password to securely link the file):
```
sudo cp pet /usr/local/bin/
```
Now you're completely set... LAUNCH!!!!

```
pet
```

## 🪟 Windows Installation (via WSL)
If you are on Windows, you can run Puchi inside a native Linux environment using the Windows Subsystem for Linux (WSL).

### Step 1: Open WSL
Open your preferred terminal app (like Windows Terminal or PowerShell) and boot into your Linux system environment:

```
wsl
```
### Step 2: Environment Setup
If you don't have a Linux distribution set up inside WSL yet, install a clean standard distribution (like Ubuntu or Debian) from the Microsoft Store or via command line, then access its terminal shell.

### Step 3: Global Installation
Once inside your WSL Linux terminal prompt, simply follow the standard Linux Installation instructions above to update your packages, build the code, and link the binary!

Once copied to /usr/local/bin/, typing pet inside your WSL environment will instantly wake your buddy up.

