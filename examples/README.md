## `lr0` examples

```sh
$ make all
...

$ .bin/01-calc-tiny "3+8-5" "3+ 8-5" "3+8*5"
0> 3+8-5        => 6
1> 3+ 8-5       => Error: unexpected input: expected int: parse error near ⟪3+⟫⏵⟪␠8-5⟫
2> 3+8*5        => Error: unexpected input: expected "+" or "-": parse error near ⟪3+8⟫⏵⟪*5⟫

$ .bin/02-calc "42* 23+17" "42*(23+17)" "3+8*)"
0> 42* 23+17    => 983
1> 42*(23+17)   => 1680
2> 3+8*)        => Error: unexpected input: expected int or "(": parse error near ⟪3+8*⟫⏵⟪)⟫

$ make clean
```
