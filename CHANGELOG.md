# Changelog

## v0.1.1

### Fixed

- SetPull: fix IOC MMIO write; write-enable mask was shifted twice (`mask<<16`)
  and lower-bit merge used the wrong mask, so pull settings never took effect.
