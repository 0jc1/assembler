main:
    lw t0, 4(sp)
    li a0 50        # first argument
    li a1 7         # second argument

    call add_numbers # call function
    li a7 1       # syscall: print integer
    ecall
    
    li a7 10
    ecall
    
add_numbers:
    add a0 a0 a1 # a0 = a0 + a1
    ret