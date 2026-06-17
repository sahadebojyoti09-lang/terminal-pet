# 🐾 Puchi - The Ultimate Terminal Pet

A blazing fast, ultra-smooth virtual terminal pet built with **Go** and the **Bubble Tea** TUI framework. Puchi lives inside your terminal, blinks, gets hungry, and dispenses random wisdom using `fortune`.

---

## ✨ Features

* **Zero Flicker:** Uses double-buffering terminal rendering for buttery-smooth animations.
* **Low Footprint:** Compiles into a single native binary with near 0% CPU idle usage.
* **Fortune-Powered:** Pressing a key makes Puchi speak short quotes directly from your system's `fortune` database.
* **Persistent Speech:** Text freezes on screen so you can read fortunes completely at your own pace.

---

## 🚀 Quick Start

### 1. Prerequisites

Make sure you have `go` and `fortune` installed on your system

### 2. Close & Run

Clone the repo
```
git clone https://github.com/sahadebojyoti09-lang/terminal-pet.git
```
Move into the pet directory
```
cd terminal-pet
```
Build it
```
go build -o pet main.go
```
Move to your local binaries path (it will ask for sudo password but it is nothing to worry about)
```
sudo cp pet /usr/local/bin/
```
Now you're set...
LAUNCH!!!!
```
pet
```


