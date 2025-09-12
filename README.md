<h4>
This project is currently EXPERIMENTAL! <br />
Please check if your optimized pack works the same as your original pack does. <br />
Please <a href="https://github.com/TBroz15/TesserPack/issues">report any issues</a> you've encountered.
</h4>


<div align="center">
    <img src="./.github/assets/tesserpack.svg" width="200" />
    <h6><code>logo made in excalidraw... pls dont judge me lol</code></h6>
    <h1>
        TesserPack
    </h1>
    <h5>
        Optimize Minecraft Resouce &amp; Behavior Packs at the speed of Ice Boats!<br>
        (get it? because its fast to travel)
    </h5>
</div>

**TesserPack** is a build tool for Minecraft that compiles and optimize any Minecraft pack to make it more easier for people to download. TesserPack can compress large packs such as [bedrock-samples](https://github.com/Mojang/bedrock-samples) by atleast twice smaller the size![^1]

## But why?

<p>There is a problem with big Minecraft packs. Where assets, scripts (in Bedrock Edition), and JSON files aren't well optimized in downloading for people who have internet issues. <b>TesserPack</b> fixes this by compressing packs as much as possible, reducing the bandwidth required to download and get micro-optimizations <i>vs. the unoptimized version of the packs.</i></p>

<details>
<summary><b>Why the name "TesserPack?"</b></summary>
<p>I simply combined "tesseract" and "pack" together. I put "tesseract" because it is a cool shape. I put "pack" because the main function is to optimize Minecraft packs.</p>
</details>

## How to Use it?

### All you need is...

- Just a transistor! (the more transistor that a computer has, the better performance and less likely to crash)

### Installation?

As of now, TesserPack is available to download in the [Releases page](https://github.com/TBroz15/TesserPack/releases/latest). The Go package will be available if TesserPack is near to complete as a project or there is nothing to add more features to TesserPack.

#### For Windows Systems,

TesserPack for Windows is ready to use and prebuilt along with required dependencies. Since it is portable, just run the executable in the terminal if you are on the same directory or if you make it globally available to your terminal.

#### For Linux Systems,

TesserPack for Linux includes the prebuilt executable itself only and not the required dependencies. You must atleast have libvips, mozjpeg, libspng, libimagequant, and highway installed, in order for TesserPack to work. You can use this [Libvips installation guide](https://github.com/TBroz15/TesserPack/blob/main/CONTRIBUTING.md#linux-setup) for reference.

### Compile through CLI

This is an example on how to optimize your pack.

```bash
./tesserpack --in ./cool-and-epic-pack/ --out ./the-swag-dir/
```

The `--in` flag defines the directory path of the pack is supposed to be.

The `--out` flag is where the output of the compiled pack will appear.

If your pack's name is `cool-and-epic-pack`, then your optimized pack will appear at `the-swag-dir` with the name `cool-and-epic-pack-optimized`.

## TesserPack vs. SuitcaseJS

TesserPack is a full Golang rewrite and the successor of [SuitcaseJS](https://github.com/TBroz15/SuitcaseJS), with atleast the same features and purpose. But it is further improved from what SuitcaseJS was missing, such as:

- Better compilation performance[^2] [^3]
- Less memory usage (up to 10x)[^2] [^4]
- Small fixes & features:
    - Copy the original image if the "optimized" version has larger size
    - Can run multiple instances [^5]
    - Recompile already compiled packs (zip/mcpack)
    - [^2]
- JSON5 support
- User-friendly GUI (coming soon)
- JS & TS minification (also coming soon)

[^1]: Comparison of the uncompiled & unzipped pack (~200MB) to an optimized & zipped pack (~100MB).
[^2]: That's because it isn't written in Javascript. 
[^3]: TesserPack was designed to handle large Minecraft packs like [bedrock-samples](https://github.com/Mojang/bedrock-samples) without some overhead.
[^4]: I've found that SuitcaseJS is using ~1.2GB vs. TesserPack with ~100MB in their peak RAM usage when compiling [bedrock-samples](https://github.com/Mojang/bedrock-samples). dat so insane vro ☠️☠️☠️ -TuxeBro
[^5]: Running multiple instances of SuitcaseJS will break since it works in a hard coded directory (*homedir*/.suitcase/temp) at the same time. TesserPack can have multiple instances but it will warn you since it can be resource intensive.