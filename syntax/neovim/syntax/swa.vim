" Vim syntax file
" Language: swa 
" Maintainer: Pathe SENE 
" Latest Revision: 2025-11-22

if exists("b:current_syntax")
  finish
endif

syn keyword swaConditionalEnglish  if else
syn keyword swaConditionalFrench   si sinon 
syn keyword swaFunctionEnglish     print 
syn keyword swaKeywordEnglish      func struct start dialect int float string if else while print let return 
syn keyword swaKeywordFrench       fonction structure demarrer dialecte entier decimal chaine si sinon afficher variable retourner 
syn keyword swaRepeatEnglish       while 
syn keyword swaRepeatFrench        tant que
syn keyword swaTypeEnglish         int float string 
syn keyword swaTypeFrench          entier decimal chaine 
syn region swaComment              start="//" skip="\\$" end="$" contains=@swaTodo

hi def link swaKeywordEnglish     Keyword
hi def link swaKeywordFrench      Keyword
hi def link swaConditionalEnglish Conditional
hi def link swaConditionalFrench Conditional
hi def link swaRepeatEnglish      Repeat
hi def link swaRepeatFrench       Repeat
hi def link swaTypeEnglish        Type
hi def link swaTypeFrench         Type
hi def link swaComment            Comment
hi def link swaFunctionEnglish    Function

let b:current_syntax = "swa"
