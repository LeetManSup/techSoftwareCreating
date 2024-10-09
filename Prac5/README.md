# Практическая работа №5. Задачи для практической работы на языке Go

## 1. Задачи на линейное программирование (без условных операторов и циклов)

### Задание 1.1.
Напишите программу, которая принимает целое число и вычисляет сумму его цифр.
```go
package main

import (
	"fmt"
	"os"
)

func main() {
	var number, sum int
	sum = 0

	fmt.Fscan(os.Stdin, &number)

	sum += number % 10
	number /= 10
	sum += number % 10
	number /= 10
	sum += number % 10
	number /= 10
	sum += number % 10

	fmt.Println(sum)
}
```

### Задание 1.2.
Напишите программу, которая преобразует температуру из градусов Цельсия в Фаренгейты и обратно. 
```go
package main

import (
	"fmt"
	"os"
)

func main() {
	var c float32
	var f float32

	fmt.Fscan(os.Stdin, &c)
	fmt.Fscan(os.Stdin, &f)

	fmt.Println(c, "град. C = ", convertToF(c), "град. F")
	fmt.Println(f, "град. F = ", convertToC(f), "град. C")

}

func convertToF(value float32) float32 {
	return 9.0/5.0*value + 32
}

func convertToC(value float32) float32 {
	return 5.0 / 9.0 * (value - 32)
}
```

### Задание 1.3.
Напишите программу, которая принимает массив чисел и возвращает новый массив, где каждое число удвоено.
```go
package main

import "fmt"

func main() {
	massive := [4]int{1, 2, 3, 4}

	massive[0] *= 2
	massive[1] *= 2
	massive[2] *= 2
	massive[3] *= 2

	fmt.Print(massive)
}
```

### Задание 1.4.
Напишите программу, которая принимает несколько строк и объединяет их в одну строку через пробел.
```go
package main

import (
	"fmt"
	"os"
)

func main() {
	var str1, str2 string

	fmt.Fscan(os.Stdin, &str1)
	fmt.Fscan(os.Stdin, &str2)

	result := str1 + " " + str2

	fmt.Println(result)
}
```

### Задание 1.5.
Напишите программу, которая вычисляет расстояние между двумя точками в 2D пространстве.
```go
package main

import (
	"fmt"
	"math"
	"os"
)

func main() {
	var x1, x2, y1, y2 float64

	fmt.Printf("Введите x1: ")
	fmt.Fscan(os.Stdin, &x1)

	fmt.Printf("Введите x2: ")
	fmt.Fscan(os.Stdin, &y1)

	fmt.Printf("Введите y1: ")
	fmt.Fscan(os.Stdin, &x2)

	fmt.Printf("Введите y2: ")
	fmt.Fscan(os.Stdin, &y2)

	result := math.Sqrt(math.Pow(x2-x1, 2) + math.Pow(y2-y1, 2))

	fmt.Println(result)
}
```

## 2. Задачи с условным оператором

### Задание 2.1.
Напишите программу, которая проверяет, является ли введенное число четным или нечетным.
```go
package main

import (
	"fmt"
	"os"
)

func main() {
	var number int
	fmt.Fscan(os.Stdin, &number)

	if number%2 == 1 {
		fmt.Println("Нечётное")
	} else {
		fmt.Println("Чётное")
	}
}
```

### Задание 2.2.
Напишите программу, которая проверяет, является ли введенный год високосным.
```go
package main

import (
	"fmt"
	"os"
)

func main() {
	var year int
	fmt.Fscan(os.Stdin, &year)

	if year%4 == 0 {
		fmt.Println("Високосный")
	} else {
		fmt.Println("Не високосный")
	}
}
```

### Задание 2.3.
Напишите программу, которая принимает три числа и выводит наибольшее из них.
```go
package main

import (
	"fmt"
	"os"
)

func main() {
	var num1, num2, num3 int

	fmt.Fscan(os.Stdin, &num1)
	fmt.Fscan(os.Stdin, &num2)
	fmt.Fscan(os.Stdin, &num3)

	if num1 > num2 {
		if num1 > num3 {
			fmt.Println(num1)
		} else {
			fmt.Println(num3)
		}
	} else if num3 < num2 {
		fmt.Println(num2)
	} else {
		fmt.Println(num3)
	}
}
```

### Задание 2.4.
Напишите программу, которая принимает возраст человека и выводит, к какой возрастной группе он относится (ребенок, подросток, взрослый, пожилой. В комментариях указать возрастные рамки).
```go
package main

import (
	"fmt"
	"os"
)

func main() {
	var age int

	fmt.Fscan(os.Stdin, &age)

	/*
	Взрослый 20+
	Подросток 13-19
	Ребёнок 0-13 
	*/

	if age > 19 {
		fmt.Println("Взрослый")
	} else if age < 13 {
		fmt.Println("Подросток")
	} else {
		fmt.Println("Ребёнок")
	}
}
```

### Задание 2.5.
Напишите программу, которая проверяет, делится ли число одновременно на 3 и 5.
```go
package main

import (
	"fmt"
	"os"
)

func main() {
	var number int

	fmt.Fscan(os.Stdin, &number)

	if number%3 == 0 && number%5 == 0 {
		fmt.Println("Делится")
	} else {
		fmt.Println("Не делится")
	}
}
```

## 3. Задачи на циклы
### Задание 3.1.
Напишите программу, которая вычисляет факториал числа.
```go
package main

import (
	"fmt"
	"os"
)

func factorial(number int) int {
	if number == 0 || number == 1 {
		return 1
	}

	return number * factorial(number-1)
}

func main() {
	var number int
	fmt.Fscan(os.Stdin, &number)

	fmt.Println(factorial(number))
}
```

### Задание 3.2.
Напишите программу, которая выводит первые "n" чисел Фибоначчи.
```go
package main

import (
	"fmt"
	"os"
)

func fibonachi(number int) int {
	if number == 0 || number == 1 {
		return number
	}

	return fibonachi(number-1) + fibonachi(number-2)
}

func main() {
	var number int
	fmt.Fscan(os.Stdin, &number)

	for i := 0; i <= number; i++ {
		fmt.Print(fibonachi(i), " ")
	}
}
```

### Задание 3.3.
Напишите программу, которая переворачивает массив чисел.
```go
package main

import (
	"fmt"
	"os"
)

func main() {
	var massive_size int

	fmt.Println("Введите размер массива")
	fmt.Fscan(os.Stdin, &massive_size)

	var massive = make([]int, massive_size)

	fmt.Println("Заполните массив")

	for i := 0; i < massive_size; i++ {
		fmt.Fscan(os.Stdin, &massive[i])
	}

	for i := 0; i < massive_size/2; i++ {
		temp := massive[i]
		massive[i] = massive[massive_size-i-1]
		massive[massive_size-i-1] = temp
	}

	fmt.Println(massive)
}
```

### Задание 3.4.
Напишите программу, которая выводит все простые числа до заданного числа.
```go
package main

import (
	"fmt"
	"os"
)

func main() {
	var num int

	fmt.Print("Введите число: ")
	fmt.Fscan(os.Stdin, &num)

	for i := 2; i <= num; i++ {
		for j := 2; j < i; j++ {
			if i%j == 0 {
				goto exit
			}
		}
		{
			fmt.Print(i, " ")
		}
	exit:
	}
}

```

### Задание 3.5.
Напишите программу, которая вычисляет сумму всех чисел в массиве.
```go
package main

import (
	"fmt"
	"os"
)

func main() {
	var massiveSize int
	var result = 0

	fmt.Print("Введите размер массива: ")
	fmt.Fscan(os.Stdin, &massiveSize)

	var massive = make([]int, massiveSize)

	fmt.Println("Заполните массив: ")

	for i := 0; i < massiveSize; i++ {
		fmt.Fscan(os.Stdin, &massive[i])
	}

	for _, value := range massive {
		result += value
	}

	fmt.Print("Сумма чисел в массиве: ", result)
}

```