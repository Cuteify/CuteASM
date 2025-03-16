section .data; 
section .text; 
; ==============================
; Function:test.hiMyLang2
test.hiMyLang2:
    PUSH ESP
    MOV EBP, ESP
    SUB ESP, 8
    MOV EBX, 
    ADD EBX, 3
    MOV EAX, EBX
    CMP EAX, 6666
    if_1:
        ADD esp, 16
        POP ebp
    end_if_1:
        MOV DWORD[ebp-8], 123
        CMP 123, EAX
    if_2:
        MOV DWORD[ebp-8], 9
    else_if_2:
        MOV DWORD[ebp-8], 10
    end_if_2:
        ADD esp, 16
        POP ebp

; Function End:test.hiMyLang2
; ==============================

; ==============================
; Function:test.hiFn2
test.hiFn2:
    PUSH ESP
    MOV EBP, ESP
    SUB ESP, 0
    PUSH ebp
    MOV ebp, esp
    SUB esp, 16
    MOV , 9
    MOV , 78
    CALL 
    MOV , 5
    MOV , 6
    MOV , 1
    if_3:
        MOV , 0
    else_if_3:
        MOV , 10
    end_if_3:
        CMP EAX, 0
    if_4:
        MOV , 9
    else_if_4:
        ADD esp, 16
        POP ebp
    end_if_4:
        CMP EAX, 0
    if_5:
        MOV , 9
    end_if_5:
        ADD esp, 16
        POP ebp

; Function End:test.hiFn2
; ==============================

; ==============================
; Function:test.print0
test.print0:
    PUSH ESP
    MOV EBP, ESP
    SUB ESP, 0
    PUSH ebp
    MOV ebp, esp
    SUB esp, 4
    PUSH 
    CALL 
    PUSH 0
    PUSH 0
    PUSH 
    PUSH 
    PUSH EAX
    CALL 
    XOR EAX, EAX
    ADD esp, 4
    POP ebp

; Function End:test.print0
; ==============================

; ==============================
; Function:test.main0
test.main0:
    PUSH ESP
    MOV EBP, ESP
    SUB ESP, 0
    PUSH ebp
    MOV ebp, esp
    SUB esp, 12
    MOV , 1
    MOV , 100
    CALL 
    CALL 
    ADD esp, 12
    POP ebp

; Function End:test.main0
; ==============================

; ==============================
; Function:main
main:
    PUSH ESP
    MOV EBP, ESP
    SUB ESP, 0
    CALL 

; Function End:main
; ==============================

