package decoder

import (
	"errors"
	"io"
	"strings"
	"testing"
	"time"

	decoderv2 "github.com/y1w5/stream/db/decoder/v2"
)

var xmlSample = `<mediawiki xmlns="http://www.mediawiki.org/xml/export-0.10/" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xsi:schemaLocation="http://www.mediawiki.org/xml/export-0.10/ http://www.mediawiki.org/xml/export-0.10.xsd" version="0.10" xml:lang="en">
  <siteinfo>
    <sitename>Wikipedia</sitename>
    <dbname>enwiki</dbname>
    <base>https://en.wikipedia.org/wiki/Main_Page</base>
    <generator>MediaWiki 1.42.0-wmf.1</generator>
    <case>first-letter</case>
    <namespaces>
      <namespace key="-2" case="first-letter">Media</namespace>
      <namespace key="-1" case="first-letter">Special</namespace>
      <namespace key="0" case="first-letter" />
      <namespace key="1" case="first-letter">Talk</namespace>
      <namespace key="2" case="first-letter">User</namespace>
      <namespace key="3" case="first-letter">User talk</namespace>
      <namespace key="4" case="first-letter">Wikipedia</namespace>
      <namespace key="5" case="first-letter">Wikipedia talk</namespace>
      <namespace key="6" case="first-letter">File</namespace>
      <namespace key="7" case="first-letter">File talk</namespace>
      <namespace key="8" case="first-letter">MediaWiki</namespace>
      <namespace key="9" case="first-letter">MediaWiki talk</namespace>
      <namespace key="10" case="first-letter">Template</namespace>
      <namespace key="11" case="first-letter">Template talk</namespace>
      <namespace key="12" case="first-letter">Help</namespace>
      <namespace key="13" case="first-letter">Help talk</namespace>
      <namespace key="14" case="first-letter">Category</namespace>
      <namespace key="15" case="first-letter">Category talk</namespace>
      <namespace key="100" case="first-letter">Portal</namespace>
      <namespace key="101" case="first-letter">Portal talk</namespace>
      <namespace key="118" case="first-letter">Draft</namespace>
      <namespace key="119" case="first-letter">Draft talk</namespace>
      <namespace key="710" case="first-letter">TimedText</namespace>
      <namespace key="711" case="first-letter">TimedText talk</namespace>
      <namespace key="828" case="first-letter">Module</namespace>
      <namespace key="829" case="first-letter">Module talk</namespace>
      <namespace key="2300" case="case-sensitive">Gadget</namespace>
      <namespace key="2301" case="case-sensitive">Gadget talk</namespace>
      <namespace key="2302" case="case-sensitive">Gadget definition</namespace>
      <namespace key="2303" case="case-sensitive">Gadget definition talk</namespace>
    </namespaces>
  </siteinfo>
  <page>
    <title>AccessibleComputing</title>
    <ns>0</ns>
    <id>10</id>
    <redirect title="Computer accessibility" />
    <revision>
      <id>1002250816</id>
      <parentid>854851586</parentid>
      <timestamp>2021-01-23T15:15:01Z</timestamp>
      <contributor>
        <username>Elli</username>
        <id>20842734</id>
      </contributor>
      <minor />
      <comment>shel</comment>
      <model>wikitext</model>
      <format>text/x-wiki</format>
      <text bytes="111" xml:space="preserve">#REDIRECT [[Computer accessibility]]

{{rcat shell|
{{R from move}}
{{R from CamelCase}}
{{R unprintworthy}}
}}</text>
      <sha1>kmysdltgexdwkv2xsml3j44jb56dxvn</sha1>
    </revision>
  </page>
  <page>
    <title>Anarchism</title>
    <ns>0</ns>
    <id>12</id>
    <revision>
      <id>1178630344</id>
      <parentid>1178595022</parentid>
      <timestamp>2023-10-04T21:56:10Z</timestamp>
      <contributor>
        <username>Cinadon36</username>
        <id>29281028</id>
      </contributor>
      <comment>/* Key issues */ removing recently added sentence. Because of ?RS and not adding anything substantial.</comment>
      <model>wikitext</model>
      <format>text/x-wiki</format>
      <text bytes="3091" xml:space="preserve">{{short description|Political philosophy and movement}}
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

'''Anarchism''' is a [[political philosophy]] and [[Political movement|movement]] that is skeptical of all justifications for [[authority]] and seeks to abolish the [[institutions]] it claims maintain unnecessary [[coercion]] and [[Social hierarchy|hierarchy]], typically including [[Nation state|nation-states]],{{sfn|Suissa|2019b|ps=: &quot;...as many anarchists have stressed, it is not government as such that they find objectionable, but the hierarchical forms of government associated with the nation state.&quot;}} and [[capitalism]]. Anarchism advocates for the replacement of the state with [[Stateless society|stateless societies]] and [[Voluntary association|voluntary]] [[Free association (communism and anarchism)|free associations]]. As a historically [[left-wing]] movement, this reading of anarchism is placed on the [[Far-left politics|farthest left]] of the [[political spectrum]], usually described as the [[libertarian]] wing of the [[socialist movement]] ([[libertarian socialism]]).

[[Humans]] have lived in [[society|societies]] without formal hierarchies long before the establishment of states, [[realm]]s, or [[empire]]s. With the rise of organised hierarchical bodies, [[scepticism]] toward authority also rose. Although traces of anarchist ideas are found all throughout history, modern anarchism emerged from the [[Age of Enlightenment|Enlightenment]]. During the latter half of the 19th and the first decades of the 20th century, the anarchist movement flourished in most parts of the world and had a significant role in workers' struggles for [[emancipation]]. [[Anarchist schools of thought|Various anarchist schools of thought]] formed during this period. Anarchists have taken part in several revolutions, most notably in the [[Paris Commune]], the [[Russian Civil War]] and the [[Spanish Civil War]], whose end marked the end of the [[classical era of anarchism]]. In the last decades of the 20th and into the 21st century, the anarchist movement has been resurgent once more, growing in popularity and influence within [[anti-capitalist]], [[anti-war]] and [[anti-globalisation]] movements.

Anarchists employ [[diversity of tactics|diverse approaches]], which may be generally divided into revolutionary and [[evolutionary strategies]]; there is significant overlap between the two. Evolutionary methods try to simulate what an anarchist society might be like, but revolutionary tactics, which have historically taken a violent turn, aim to overthrow authority and the state. Many facets of human civilization have been influenced by anarchist theory, critique, and [[Praxis (process)|praxis]].
{{toc limit|3}}</text>
      <sha1>3ow6ges2gq0tv5oojdf1s6du0wbdd4t</sha1>
    </revision>
  </page>
</mediawiki>`

var pages = []Page{
	{
		UpdatedAt: mustParseTime(time.RFC3339, "2021-01-23T15:15:01Z"),
		Title:     "AccessibleComputing",
		Text: `#REDIRECT [[Computer accessibility]]

{{rcat shell|
{{R from move}}
{{R from CamelCase}}
{{R unprintworthy}}
}}`,
	},
	{
		UpdatedAt: mustParseTime(time.RFC3339, "2023-10-04T21:56:10Z"),
		Title:     "Anarchism",
		Text: `{{short description|Political philosophy and movement}}
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
{{toc limit|3}}`,
	},
}

func TestDecoder(t *testing.T) {
	r := strings.NewReader(xmlSample)
	decoder, err := New(r)
	if err != nil {
		t.Fatalf("fail to create decoder: %v", err)
	}
	for i := 0; decoder.Next(); i++ {
		var p Page
		if err := decoder.Scan(&p); err != nil {
			t.Fatalf("fail to scan page: %v", err)
		}
		if p.Text != pages[i].Text {
			t.Logf("expects.Text=%q", pages[i].Text)
			t.Logf("got.Text=%q", p.Text)
			t.Fatalf("unexpected page content")
		}
	}
	if err := decoder.Err(); err != nil && !errors.Is(err, io.EOF) {
		t.Fatalf("fail to read next page: %v", err)
	}
}

func TestDecoderV2(t *testing.T) {
	r := strings.NewReader(xmlSample)
	d, err := decoderv2.New(r)
	if err != nil {
		t.Fatalf("fail to create decoder: %v", err)
	}
	for i := 0; d.Next(); i++ {
		var p decoderv2.Page
		if err := d.Scan(&p); err != nil {
			t.Fatalf("fail to scan page: %v", err)
		}
		if p != decoderv2.Page(pages[i]) {
			t.Logf("expects=%v", pages[i])
			t.Logf("got=%v", p)
			t.Fatalf("unexpected page content")
		}
	}
	if err := d.Err(); err != nil && !errors.Is(err, io.EOF) {
		t.Fatalf("fail to read next page: %v", err)
	}
}

func mustParseTime(layout, value string) time.Time {
	t, err := time.Parse(layout, value)
	if err != nil {
		panic(err)
	}
	return t
}
