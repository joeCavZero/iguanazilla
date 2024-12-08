.data
macaco: .asciiz "tome   "
sagui: .asciiz "hi"#HAPPY
chipanze: .word 20,30,40,50
bonobo: .word 10, 11
mico: .word 1 , 2, 3 , 4,      5
orango: .byte 30

#MACACOOOOO
.text
loop:
    ADDD           20    #meo
SUBD 10
MULD#OI