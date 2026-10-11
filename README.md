# swayom (Sway Output Manager)

This is a tool for managing your outputs in sway.

Normally you define a configuration for each of your displays in sway and
these configurations are applied at every session refresh. However it is
not possible to have different configurations based on the connected
displays.

swayom solves this issue by listening for output events and loading a
custom configuration every time an output change happens, such as plugging
in/off a display.

## Usage

An example configuration looks like this:

```
@@["MAKE MODEL ID1", "MAKE MODEL ID2"]
output 'MAKE MODEL ID1' {
  mode 1920x1080@60Hz
  position 0 0
}
output 'MAKE MODEL ID2' {
  mode 1920x1080@60Hz
  position 1920 0
}
@@
@@["MAKE MODEL ID1", "MAKE MODEL ID2", "MAKE MODEL ID3"]
output 'MAKE MODEL ID1' {
  mode 1920x1080@60Hz
  position 0 0
}
output 'MAKE MODEL ID3' {
  mode 1920x1080@60Hz
  position 1920 0
}
output 'MAKE MODEL ID2' {
  mode 1920x1080@60Hz
  position 3840 0
}
@@
```

In its current state, only the full output IDs are accepted as markers but
any configuration may be written in the block. The configuration is applied
if the outputs in the json array match with the currently active outputs.
