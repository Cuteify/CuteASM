section .data
;    message db 'Hello, World!', 0
;    messageLen equ $ - message

section .text
;global main


; ==============================
; Function:test.hiMyLang2
test.hiMyLang2:(qw hi, dw b)
    var $a, dw
    var $b, dw
    mov %e1, DW[%rbp+12]; 保存表达式左边的值
    add %e1, 3; 计算表达式的值
    mov %e0, %e1; 
    cmp %e0, 6666; 比较表达式的值
    jnl end_if_1; 判断后跳转到目标
    if_1:
    add %rsp, 16; 还原栈指针
    pop %rbp; 跳转到函数返回部分
    ret

    end_if_1:
    mov $b, 123; 设置变量
    cmp 123, %e0; 比较表达式的值
    jnl else_if_2; 判断后跳转到目标
    if_2:
    mov  $b, 9; 设置变量
    else_if_2:
    mov  $b, 10; 设置变量
    end_if_2:
    add %rsp, 16; 还原栈指针
    pop %rbp; 跳转到函数返回部分
    ret

; ======函数完毕=======

; ==============================
; Function:test.hiFn2
test.hiFn2:()
    push %rbp; 函数基指针入栈
    mov %rbp, %rsp; 设置基指针
    sub %rsp, 16; 调整栈指针
    mov DW[%rsp+8], 9; 设置函数参数
    mov DW[%rsp+4], 78; 设置函数参数
    call test.hiMyLang2; 调用函数
    mov  DW[%rbp-4], 5; 设置变量
    mov  DW[%rbp-8], 6; 设置变量
    ;mov else_if_3, 1; 
    if_3:
    mov  DW[%rbp-8], 0; 设置变量
    else_if_3:
    mov  DW[%rbp-8], 10; 设置变量
    end_if_3:
    cmp %e0, 0; 比较表达式的值
    jnl else_if_4; 判断后跳转到目标
    if_4:
    mov  DW[%rbp-8], 9; 设置变量
    else_if_4:
    add %rsp, 16; 还原栈指针
    pop %rbp; 跳转到函数返回部分
    ret

    end_if_4:
    cmp %e0, 0; 比较表达式的值
    jnl end_if_5; 判断后跳转到目标
    if_5:
    mov  DW[%rbp-8], 9; 设置变量
    end_if_5:
    add %rsp, 16; 还原栈指针
    pop %rbp; 弹出函数基指针
    ret

; ======函数完毕=======

; ==============================
; Function:test.print0
test.print0:()
    push %rbp; 函数基指针入栈
    mov %rbp, %rsp; 设置基指针
    sub %rsp, 4; 调整栈指针
    push -11; 设置函数参数
    extern  GetStdHandle@1; 外部函数
    call  GetStdHandle@1; 调用外部函数

extern WriteFile

    ; 准备WriteConsoleW的参数
    push 0      ; 第五个参数：不保留额外
    push 0          ; 第四个参数：缓冲区，用于接收实际写入的字节数
    ;push messageLen ; 第三个参数：字符串的长度
    ;push message ; 第二个参数：指向要写入的字符串
    push %e0      ; 第一个参数：句柄
    
    ; 调用WriteFileW
    call WriteFile

    xor %e0, %e0
    
    add %rsp, 4; 还原栈指针
    pop %rbp; 弹出函数基指针
    ret

; ======函数完毕=======

; ==============================
; Function:test.main0
test.main0:()
    push %rbp; 函数基指针入栈
    mov %rbp, %rsp; 设置基指针
    sub %rsp, 12; 调整栈指针
    mov DW[%rsp+12], 1; 设置函数参数
    mov DW[%rsp+8], 100; 设置函数参数
    call test.hiFn2; 调用函数
    call test.print0; 调用函数
    add %rsp, 12; 还原栈指针
    pop %rbp; 弹出函数基指针
    ret

; ======函数完毕=======


main:()
call test.main0()
;PRINT_STRING "MyLang First Finish!"
ret

