<p align="center">
  <img align="center" width="100%" height="100%" src="interpreter.png">
</p>

# What is SkyLine 

SkyLine or CSC ( Cyber Security Core ) is a new interpreted programming language dedicated to just cyber security and mathematics. SkyLine plans to change the way the cyber security world operates by provided a much more slick and modern interface for security researches or cyber security solutions developers.

# Why just cybersecurity 

Cyber Security is a major part of the developers career, after all the whole reason SkyLine exists is due to CyberSecurity interests. However it was felt amongst the developers of the language that there was not a language out there that was really modern and well rounded for anything cyber security. There was just general purpose scripting languages liek Python or Perl which had libraries dedicated to cyber security which do not always work in the cases that they are developed. SkyLine wants to change that by standardizing both offensive and defensive cyber security.

# Operating system support 

Currenty SkyLine has a very small and limited interpreter however despite it not being OS dependant the current error system and banner system is a little wack if you are on windows. SkyLine has only been tested on windows as its main development purpose was for linux.

# Documentation 

SkyLine has many ways of allowing the user to full customization of the style within the language, due to this abstract design SkyLine has some pretty interesting ways of writing code. Below are examples and basic documentation on the SkyLine programming language.

> Declaring variables 

* About: SkyLine has many ways of defining variables, types include float, string, bool and integer as of 0.0.2. Def keywords consist of `allow, let, cause`

here is a few ways you can declare variables in SkyLine 

```rs
allow Variable = "data"
let   Variable = "data"
cause Variable = "data"
```

cause, let and allow all do the same thing they are just different keywords

> Notes 

* Notes in SkyLine are not fully supported, as whitespace within the comment line is not fully yet finished. You can define a note like so 

```rs
!!anote
#anote
//anote
```

if you create a note that has spaces, like so 

```rs
!! hello there, ext, function name, name, variable
```

SkyLine will tell you there is an error in your statement such as `no parse prefix for !! found`. This is because comments are not fully supported

> Writing functions 

* Note: Functions in skyline can be developed and written in all sorts of ways but there are three main keywords which are `fn, Func and function`

Defining a function is like so

```rs

```
