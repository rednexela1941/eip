#!/usr/bin/env perl
use strict;
use warnings;

use v5.32;

my @types = qw(
  BOOL
  SINT
  INT
  DINT
  LINT

  USINT
  UINT
  UDINT
  ULINT

  BYTE
  WORD
  DWORD
  LWORD
);

foreach (@types)
{
    say <<HERE;
func New${_}Param(name string, ptr *cip.$_) *AssemblyParam {
	return _NewDefaultParam(name, $_, ptr)
}
HERE
}

foreach (@types)
{
    say <<HERE;
func (self *AssemblyInstance) Add${_}Param(name string, ptr *cip.$_) *param.AssemblyParam {
  p := param.New${_}Param(name, ptr)
  self.AddParam(p)
  return p
}
HERE

}
