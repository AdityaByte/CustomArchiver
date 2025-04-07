# 📦 AdZip - Your Minimalist Archiver  
*A no-nonsense tool to bundle files together. No fancy compression, just simple archiving.*

## 🔧 How to Use

### Build It:
```bash
go build -o adzip main.go
```

### 1. Archive Files  
```bash
./adzip -archive "file1.txt image.png" -o bundle 
```
→ Creates bundle.adzip containing both files.

### 2. UnArchive Files
```bash 
./adzip -unarchive bundle.adzip
```
→ Extracts files into a storage/ folder.

## 🔧 How to Use
Stores files in a simple structured format:
```
Filename  
Size  
::END-METADATA::  
[File Data]  
::END-FILE::  
```
Extracted files retain their original names.

## ❓ Why?
For when you just need a quick way to bundle files without compression

Useful for custom backups or simple data packaging

## 🚧 Limitations
No compression (files are stored as-is)

No folder support (flat files only)