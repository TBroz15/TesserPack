## TODO (for TuxeBro)

### Before v1.0

- [x] Implement main features (such as JSON, assets, and zip compression)
- [x] JSON5 support (pack devs will like this)
- [x] Proper Linux Support (libvips is the problem)
- [x] Caching
- [ ] Config
- [x] (Not so) Proper error checking
- [ ] v0 GUI w/ wails
- [ ] Full Refactor for v1

### After v1.0

- [ ] v1 GUI (when TesserPack-GUI is considered to be complete)
- [ ] Minify JS & TS through ESBuild (TBD by v1.1)
- [ ] Graceful Shutdown
- [ ] Do some optimizations and memory leak fix (ive learnt a lesson to not optimize until all the program features is ready)
    - [ ] Thread Pool for Multithreading (also add option for unli-Goroutines)
    - [ ] Multithreaded assets processing
    - [ ] Streamlike way of compiling
- [ ] Regolith Integration ([#1](https://github.com/TBroz15/TesserPack/issues/1))
- [ ] Audio Optimization?
- [ ] ARM64 Support
- [ ] TesserPack as Go package (kinda possible, just the libvips is the problem)
