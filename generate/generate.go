package generate

import (
	"fmt"
	"io"

	"github.com/dave/jennifer/jen"
)

type Enum struct {
	Name       string
	Values     []Value
	Sourcename string
}

type Value struct {
	Symbol string
	Value  string
	Pretty string
}

func Generate(packageName string, enums []Enum, w io.Writer) error {
	f := jen.NewFile(packageName)

	f.HeaderComment("Code generated by enumerator. DO NOT EDIT.")

	for _, enum := range enums {
		name := enum.Name
		values := enum.Values
		sourceName := enum.Sourcename

		errorMessage := fmt.Sprintf("invalid value for type '%s': '%%s'", name)

		f.Type().
			Id(name).
			Struct(jen.Id("v").Id(sourceName))

		f.Var().
			DefsFunc(func(g *jen.Group) {
				for _, v := range values {
					g.Id(v.Symbol).
						Op("=").
						Id(name).
						Values(jen.Id(v.Value))
				}
			})

		f.Func().
			Params(
				jen.Id("r").Id(name),
			).
			Id("Format").
			Params(
				jen.Id("f").Qual("fmt", "State"),
				jen.Id("verb").Id("rune"),
			).
			Block(
				jen.Switch(jen.Id("verb")).BlockFunc(func(g *jen.Group) {
					hasPrettyValues := false
					for _, v := range values {
						if v.Pretty != "" {
							hasPrettyValues = true
							break
						}
					}

					g.Case(jen.LitRune('s')).Block(
						jen.Id("fmt").Dot("Fprint").Call(
							jen.Id("f"), jen.Id("r").Dot("v"),
						),
					)

					g.Case(jen.LitRune('q')).Block(
						jen.Id("fmt").Dot("Fprintf").Call(
							jen.Id("f"),
							jen.Lit("%q"),
							jen.Id("r").Dot("String").Call(),
						),
					)

					if hasPrettyValues {
						g.Case(jen.LitRune('v')).Block(
							jen.Switch(jen.Id("r")).BlockFunc(func(g *jen.Group) {
								for _, v := range values {
									var val *jen.Statement
									if v.Pretty == "" {
										val = jen.Id("string").Call(jen.Id("r").Dot("v"))
									} else {
										val = jen.Lit(v.Pretty)
									}

									g.Case(jen.Id(v.Symbol)).Block(
										jen.Id("fmt").Dot("Fprint").Call(jen.Id("f"), val),
									)
								}

								g.Default().Block(
									jen.Id("fmt").Dot("Fprint").Call(jen.Id("f"), jen.Lit("")),
								)
							}),
						)
					}

					g.Default().Block(
						jen.Id("fmt").Dot("Fprint").Call(jen.Id("f"), jen.Id("r").Dot("v")),
					)
				}),
			)

		f.Func().
			Params(
				jen.Id("r").Id(name),
			).
			Id("String").
			Params().
			String().
			Block(
				jen.Return(jen.Id("string").Call(jen.Id("r").Dot("v"))),
			)

		f.Func().
			Params(
				jen.Id("r").Id(name),
			).
			Id("MarshalText").
			Params().
			Params(
				jen.Id("[]byte"),
				jen.Id("error"),
			).
			Block(
				jen.Return(
					jen.Id("[]byte").Call(jen.Id("r").Dot("v")),
					jen.Nil(),
				),
			)

		f.Func().
			Params(
				jen.Id("r").Op("*").Id(name),
			).
			Id("UnmarshalText").
			Params(
				jen.Id("__iNpUt__").Id("[]byte"),
			).
			Params(
				jen.Id("error"),
			).
			Block(
				jen.List(
					jen.Id("s"),
					jen.Id("err"),
				).Op(":=").Id("New"+name).Call(jen.Id("string").Call(jen.Id("__iNpUt__"))),

				jen.If(jen.Id("err").Op("!=").Nil()).Block(
					jen.Return(jen.Id("err")),
				),

				jen.Id("*r").Op("=").Id("s"),

				jen.Return(
					jen.Nil(),
				),
			)

		f.Func().
			Params(
				jen.Id("r").Id(name),
			).
			Id("Value").
			Params().
			Params(
				jen.Qual("database/sql/driver", "Value"),
				jen.Id("error"),
			).
			Block(
				jen.Return(
					jen.Id("r").Dot("v"),
					jen.Nil(),
				),
			)

		f.Func().
			Params(
				jen.Id("r").Op("*").Id(name),
			).
			Id("Scan").
			Params(
				jen.Id("__iNpUt__").Any(),
			).
			Params(
				jen.Id("error"),
			).
			Block(
				jen.List(
					jen.Id("s"),
					jen.Id("err"),
				).Op(":=").Id("New"+name).Call(jen.Qual("fmt", "Sprint").Call(jen.Id("__iNpUt__"))),

				jen.If(jen.Id("err").Op("!=").Nil()).Block(
					jen.Return(jen.Id("err")),
				),

				jen.Id("*r").Op("=").Id("s"),

				jen.Return(
					jen.Nil(),
				),
			)

		f.Func().
			Id("New"+name).
			Params(jen.Id("__iNpUt__").Id("string")).
			Params(
				jen.Id(name),
				jen.Id("error"),
			).
			Block(
				jen.Switch(jen.Id("__iNpUt__")).BlockFunc(func(g *jen.Group) {
					for _, v := range values {
						g.Case(jen.Id("string").Call(jen.Id(v.Value))).Block(
							jen.Return(jen.Id(v.Symbol), jen.Nil()),
						)
					}

					g.Default().Block(jen.Return(
						jen.Id(name).
							Values(),

						jen.Qual("fmt", "Errorf").Call(
							jen.Lit(errorMessage),
							jen.Id("__iNpUt__"),
						),
					))
				}),
			)
	}

	if err := f.Render(w); err != nil {
		return err
	}

	return nil
}
