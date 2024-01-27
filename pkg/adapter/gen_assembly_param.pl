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

  REAL
  LREAL
);

say "param.go Functions:";
foreach (@types)
{
    say <<HERE;
func New${_}Param(name string) *AssemblyParam {
	return NewDefaultParam(name, $_)
}
HERE
}

say "AssemblyParam Methods:";
foreach (@types)
{
    say <<HERE;
func (self *AssemblyInstance) Add${_}Param(name string) *ElementaryParam[cip.${_}] {
	ep := _NewElementaryParam[cip.${_}](name, param.${_})
	self.AddParam(ep.AssemblyParam)
	return ep
}
HERE

}

# foreach (@types) {
# 	say <<HERE;
# case ${_}Code:
# 	return ${_}Size
# HERE
# }
