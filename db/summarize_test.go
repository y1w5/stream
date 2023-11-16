package main

import "testing"

func TestSummarize(t *testing.T) {
	tests := []struct {
		name    string
		text    string
		summary string
	}{
		{
			name:    "curly",
			text:    "{{pp-semi-indef}}",
			summary: "",
		},
		{
			name: "imbricated-curly",
			text: `{{rcat shell|
{{R from move}}
{{R from CamelCase}}
{{R unprintworthy}}
}}`,
			summary: "",
		},
		{
			name: "redirect",
			text: `#REDIRECT [[Computer accessibility]]

{{rcat shell|
{{R from move}}
{{R from CamelCase}}
{{R unprintworthy}}
}}`,
			summary: "#REDIRECT [[Computer accessibility]]",
		},
		{
			name:    "header",
			text:    anarchismText,
			summary: anarchismSummary,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			summary := Summarize(tt.text)
			if summary != tt.summary {
				t.Logf("expects=%q", tt.summary)
				t.Logf("got=%q", summary)
				t.Fatalf("unexpected summary")
			}
		})
	}
}

var anarchismText = `{{short description|Political philosophy and movement}}
{{other uses}}
{{redirect2|Anarchist|Anarchists|other uses|Anarchist (disambiguation)}}
{{About|the philosophy against authority|the state of government without authority|Anarchy}}
{{pp-semi-indef}}
{{good article}}
{{use British English|date=August 2021}}
{{use dmy dates|date=August 2021}}
{{Use shortened footnotes|date=May 2023}}
{{anarchism sidebar}}
{{basic forms of government}}

'''Anarchism''' is a [[political philosophy]] and [[Political movement|movement]] that is skeptical of all justifications for [[authority]] and seeks to abolish the [[institutions]] it claims maintain unnecessary [[coercion]] and [[Social hierarchy|hierarchy]], typically including [[Nation state|nation-states]],{{sfn|Suissa|2019b|ps=: "...as many anarchists have stressed, it is not government as such that they find objectionable, but the hierarchical forms of government associated with the nation state."}} and [[capitalism]]. Anarchism advocates for the replacement of the state with [[Stateless society|stateless societies]] and [[Voluntary association|voluntary]] [[Free association (communism and anarchism)|free associations]]. As a historically [[left-wing]] movement, this reading of anarchism is placed on the [[Far-left politics|farthest left]] of the [[political spectrum]], usually described as the [[libertarian]] wing of the [[socialist movement]] ([[libertarian socialism]]).

[[Humans]] have lived in [[society|societies]] without formal hierarchies long before the establishment of states, [[realm]]s, or [[empire]]s. With the rise of organised hierarchical bodies, [[scepticism]] toward authority also rose. Although traces of anarchist ideas are found all throughout history, modern anarchism emerged from the [[Age of Enlightenment|Enlightenment]]. During the latter half of the 19th and the first decades of the 20th century, the anarchist movement flourished in most parts of the world and had a significant role in workers' struggles for [[emancipation]]. [[Anarchist schools of thought|Various anarchist schools of thought]] formed during this period. Anarchists have taken part in several revolutions, most notably in the [[Paris Commune]], the [[Russian Civil War]] and the [[Spanish Civil War]], whose end marked the end of the [[classical era of anarchism]]. In the last decades of the 20th and into the 21st century, the anarchist movement has been resurgent once more, growing in popularity and influence within [[anti-capitalist]], [[anti-war]] and [[anti-globalisation]] movements.

Anarchists employ [[diversity of tactics|diverse approaches]], which may be generally divided into revolutionary and [[evolutionary strategies]]; there is significant overlap between the two. Evolutionary methods try to simulate what an anarchist society might be like, but revolutionary tactics, which have historically taken a violent turn, aim to overthrow authority and the state. Many facets of human civilization have been influenced by anarchist theory, critique, and [[Praxis (process)|praxis]].
{{toc limit|3}}

== Etymology, terminology, and definition ==
{{main|Definition of anarchism and libertarianism}}
{{see also|Glossary of anarchism}}
[[File:WilhelmWeitling.jpg|thumb|[[Wilhelm Weitling]] is an example of a writer who added to anarchist theory without using the exact term.{{sfn|Carlson|1972|pp=22–23}}]]
The etymological origin of ''anarchism'' is from the Ancient Greek ''anarkhia'', meaning "without a ruler", composed of the prefix ''an-'' ("without") and the word ''arkhos'' ("leader" or "ruler"). The suffix ''[[-ism]]'' denotes the ideological current that favours [[anarchy]].{{sfnm|1a1=Bates|1y=2017|1p=128|2a1=Long|2y=2013|2p=217}} ''Anarchism'' appears in English from 1642 as ''anarchisme'' and ''anarchy'' from 1539; early English usages emphasised a sense of disorder.{{sfnm|1a1=Merriam-Webster|1y=2019|1loc="Anarchism"|2a1=''Oxford English Dictionary''|2y=2005|2loc="Anarchism"|3a1=Sylvan|3y=2007|3p=260}} Various factions within the [[French Revolution]] labelled their opponents as ''anarchists'', although few such accused shared many views with later anarchists. Many revolutionaries of the 19th century such as [[William Godwin]] (1756–1836) and [[Wilhelm Weitling]] (1808–1871) would contribute to the anarchist doctrines of the next generation but did not use ''anarchist'' or ''anarchism'' in describing themselves or their beliefs.{{sfn|Joll|1964|pp=27–37}}`

var anarchismSummary = `'''Anarchism''' is a [[political philosophy]] and [[Political movement|movement]] that is skeptical of all justifications for [[authority]] and seeks to abolish the [[institutions]] it claims maintain unnecessary [[coercion]] and [[Social hierarchy|hierarchy]], typically including [[Nation state|nation-states]], and [[capitalism]]. Anarchism advocates for the replacement of the state with [[Stateless society|stateless societies]] and [[Voluntary association|voluntary]] [[Free association (communism and anarchism)|free associations]]. As a historically [[left-wing]] movement, this reading of anarchism is placed on the [[Far-left politics|farthest left]] of the [[political spectrum]], usually described as the [[libertarian]] wing of the [[socialist movement]] ([[libertarian socialism]]).

[[Humans]] have lived in [[society|societies]] without formal hierarchies long before the establishment of states, [[realm]]s, or [[empire]]s. With the rise of organised hierarchical bodies, [[scepticism]] toward authority also rose. Although traces of anarchist ideas are found all throughout history, modern anarchism emerged from the [[Age of Enlightenment|Enlightenment]]. During the latter half of the 19th and the first decades of the 20th century, the anarchist movement flourished in most parts of the world and had a significant role in workers' struggles for [[emancipation]]. [[Anarchist schools of thought|Various anarchist schools of thought]] formed during this period. Anarchists have taken part in several revolutions, most notably in the [[Paris Commune]], the [[Russian Civil War]] and the [[Spanish Civil War]], whose end marked the end of the [[classical era of anarchism]]. In the last decades of the 20th and into the 21st century, the anarchist movement has been resurgent once more, growing in popularity and influence within [[anti-capitalist]], [[anti-war]] and [[anti-globalisation]] movements.

Anarchists employ [[diversity of tactics|diverse approaches]], which may be generally divided into revolutionary and [[evolutionary strategies]]; there is significant overlap between the two. Evolutionary methods try to simulate what an anarchist society might be like, but revolutionary tactics, which have historically taken a violent turn, aim to overthrow authority and the state. Many facets of human civilization have been influenced by anarchist theory, critique, and [[Praxis (process)|praxis]].`
