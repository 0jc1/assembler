.text 
    add x1, x2, x3
    addi x1, x2, 10      # x1 = x2 + 10
    lw   x1, 8(x2)       # x1 = mem[x2 + 8]
    BEQ x1, x2, label
    JAL x1, 100
label:
    add x1, x2, x3