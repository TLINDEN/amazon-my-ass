# Amazon my ass!

We europeans are not amused of Trumps trade wars, his fascistic regime
and the fact that major companies even support him. One of those
fuckers is Amazon. Leaving them is easy since we have enough
alternatives in the EU and elsewhere. Nobody needs US based morons to
buy stuff made in China anyway.

But there's one thing we need to take with us: our wish lists.

This little tool helps with this. It downloads a wishlist and saves it
to a local HTML file.

## Installation

There are multiple ways to install **amazon-my-ass**:

- Go to the [latest release page](https://github.com/tlinden/amazon-my-ass/releases/latest),
  locate the binary for your operating system and platform.
  
  Download it and put it into some directory within your `$PATH` variable.
  
- The release page also contains a tarball for every supported platform. Unpack it
  to some temporary directory, extract it and execute the following command inside:
  ```shell
  sudo make install
  ```
  
- You can also install from source. Issue the following commands in your shell:
  ```shell
  git clone https://github.com/TLINDEN/amazon-my-ass.git
  cd amazon-my-ass
  make
  sudo make install
  ```

- Or, if you have the GO toolkit installed, just install it like this:
  ```shell
  go install github.com/tlinden/amazon-my-ass@latest
  ```

If you  do not find a  binary release for your  platform, please don't
hesitate to ask me about it, I'll add it.

## Usage

* Login to your Amazon account
* Make the wishlist you want to save public
* Execute `amazon-my-ass $url` where `$url` is the URL to you public
  wishlist.
  
That's it.

## Getting help

Although I'm happy to hear from amazon-my-ass users in private email, that's the
best way for me to forget to do something.

In order to report a bug,  unexpected behavior, feature requests or to
submit    a    patch,    please    open   an    issue    on    github:
https://github.com/TLINDEN/amazon-my-ass/issues.

## Copyright and license

This software is licensed under the GNU GENERAL PUBLIC LICENSE version 3.

## Authors

T.v.Dein <tom AT vondein DOT org>

## Project homepage

https://github.com/TLINDEN/amazon-my-ass

## Copyright and License

Licensed under the GNU GENERAL PUBLIC LICENSE version 3.

## Author

T.v.Dein <tom AT vondein DOT org>
