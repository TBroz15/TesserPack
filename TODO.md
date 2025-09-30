## TODO (for TuxeBro)

### Before v1.0

- [x] Implement main features (such as JSON, assets, and zip compression)
- [x] JSON5 support (pack devs will like this)
- [x] Proper Linux Support (libvips is the problem)
- [x] Caching
- [x] Config
- [x] (Not so) Proper error checking
- [ ] v0 GUI w/ wails
- [ ] Do some optimizations and memory leak fix (ive learnt a lesson to not optimize until all the program features is ready)
    - [ ] Different Method for Multithreading
        - [ ] Thread Pool for Multithreading
        - [ ] Unbound Goroutines
    - [ ] Concurrent assets processing (remove mutex & use semaphore)
    - [ ] Metadata Caching (use hashing as fallback)
    - [ ] Streamlike way of compiling
- [ ] Github Wiki for Effective Usage 
- [ ] Add Tests
- [ ] Full Refactor for v1
    - [ ] Check if TesserPack is feasible for v1

### After v1.0

- [ ] v1 GUI
- [ ] Minify JS & TS through ESBuild (TBD by v1.1)
- [ ] Graceful Shutdown
- [ ] Regolith Integration ([#1](https://github.com/TBroz15/TesserPack/issues/1))
- [ ] Audio Optimization

### Near To Complete (kinda hard but possible)

- [ ] Mac Support
- [ ] ARM64 Support
- [ ] TesserPack as Go package (kinda possible, just the libvips is the problem)
- [ ] Full Minecraft Java Edition Support
