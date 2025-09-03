	.text
	.file	"my_module"
	.globl	main                            # -- Begin function main
	.p2align	4, 0x90
	.type	main,@function
main:                                   # @main
	.cfi_startproc
# %bb.0:                                # %func-body
	subq	$56, %rsp
	.cfi_def_cfa_offset 64
	movabsq	$8297073567416218405, %rax      # imm = 0x7325207325207325
	movq	%rax, 19(%rsp)
	movb	$0, 33(%rsp)
	movw	$25637, 31(%rsp)                # imm = 0x6425
	movl	$544417056, 27(%rsp)            # imm = 0x20732520
	movabsq	$7020093066257589614, %rax      # imm = 0x616C61732077656E
	movq	%rax, 34(%rsp)
	movabsq	$8243593322982046066, %rax      # imm = 0x7267207369207972
	movq	%rax, 42(%rsp)
	movl	$1702125925, 50(%rsp)           # imm = 0x65746165
	movw	$114, 54(%rsp)
	movl	$1851877492, 7(%rsp)            # imm = 0x6E616874
	movb	$0, 11(%rsp)
	movw	$31090, 16(%rsp)                # imm = 0x7972
	movl	$1634492787, 12(%rsp)           # imm = 0x616C6173
	movb	$0, 18(%rsp)
	movq	salary_string@GOTPCREL(%rip), %r8
	movq	salary@GOTPCREL(%rip), %r9
	leaq	19(%rsp), %rdi
	leaq	34(%rsp), %rsi
	leaq	7(%rsp), %rdx
	leaq	12(%rsp), %rcx
	xorl	%eax, %eax
	callq	printf@PLT
	movl	$10, %eax
	addq	$56, %rsp
	.cfi_def_cfa_offset 8
	retq
.Lfunc_end0:
	.size	main, .Lfunc_end0-main
	.cfi_endproc
                                        # -- End function
	.type	salary_string,@object           # @salary_string
	.data
	.globl	salary_string
salary_string:
	.asciz	"100000 EUR"
	.size	salary_string, 11

	.type	salary,@object                  # @salary
	.globl	salary
	.p2align	2, 0x0
salary:
	.long	500000                          # 0x7a120
	.size	salary, 4

	.type	new_salary,@object              # @new_salary
	.globl	new_salary
	.p2align	2, 0x0
new_salary:
	.long	1009                            # 0x3f1
	.size	new_salary, 4

	.section	".note.GNU-stack","",@progbits
