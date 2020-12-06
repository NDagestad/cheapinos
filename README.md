# Hay ! UwU

This is a cool script to find the cheapest combination of ingredients on dominos
pizzas.

The exemple data is for dominos in france but it should work everywhere as long
as they have the same system (you can replace any ingredient by any other and
the difference in price will be added to the order)

I made this for myself and if it doesn't work for you and you spent more than
the optimal amount, well to bad, send a patch ! ¯\\\_(ツ)\_/¯
(If you are seeing this on https://git.dagestad.fr/~nicolai/cheapinos the
mailling list is broken but you can send me a patch at nicolai.whathever@dagestad.fr)

# How To

```sh
go build
```

That's it (it should at least be that easy)

```sh
./cheapinos -h
```
To see how to use it

# Json files

The ingredients file should be a json array of objects with a `name: string` 
and `price: int` attributs.

The menu file should be a json array of objects with a `name: string`, 
`price: int` and `ingredients: []string` attributs


