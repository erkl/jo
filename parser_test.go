package jo

import (
	"fmt"
	"testing"
)

func ExampleParser() {
	input := []byte(`{ "foo": 10 }`)

	p := Parser{} // no initialization required
	r := 0

	for r < len(input) {
		n, ev, _ := p.Next(input[r:])
		r += n
		fmt.Printf("at %2d -> %s\n", r, ev)
	}

	ev, _ := p.End()
	fmt.Printf("at %2d -> %s\n", r, ev)

	// Output:
	// at  1 -> ObjectStart
	// at  3 -> KeyStart
	// at  7 -> KeyEnd
	// at 10 -> NumberStart
	// at 11 -> NumberEnd
	// at 13 -> ObjectEnd
	// at 13 -> Done
}

var sample = []byte(`
[
  {
    "coordinates": null,
    "favorited": false,
    "truncated": false,
    "created_at": "Wed Aug 29 17:12:58 +0000 2012",
    "id_str": "240859602684612608",
    "entities": {
      "urls": [
        {
          "expanded_url": "https://dev.twitter.com/blog/twitter-certified-products",
          "url": "https://t.co/MjJ8xAnT",
          "indices": [
            52,
            73
          ],
          "display_url": "dev.twitter.com/blog/twitter-c\u2026"
        }
      ],
      "hashtags": [

      ],
      "user_mentions": [

      ]
    },
    "in_reply_to_user_id_str": null,
    "contributors": null,
    "text": "Introducing the Twitter Certified Products Program: https://t.co/MjJ8xAnT",
    "retweet_count": 121,
    "in_reply_to_status_id_str": null,
    "id": 240859602684612608,
    "geo": null,
    "retweeted": false,
    "possibly_sensitive": false,
    "in_reply_to_user_id": null,
    "place": null,
    "user": {
      "profile_sidebar_fill_color": "DDEEF6",
      "profile_sidebar_border_color": "C0DEED",
      "profile_background_tile": false,
      "name": "Twitter API",
      "profile_image_url": "http://a0.twimg.com/profile_images/2284174872/7df3h38zabcvjylnyfe3_normal.png",
      "created_at": "Wed May 23 06:01:13 +0000 2007",
      "location": "San Francisco, CA",
      "follow_request_sent": false,
      "profile_link_color": "0084B4",
      "is_translator": false,
      "id_str": "6253282",
      "entities": {
        "url": {
          "urls": [
            {
              "expanded_url": null,
              "url": "http://dev.twitter.com",
              "indices": [
                0,
                22
              ]
            }
          ]
        },
        "description": {
          "urls": [

          ]
        }
      },
      "default_profile": true,
      "contributors_enabled": true,
      "favourites_count": 24,
      "url": "http://dev.twitter.com",
      "profile_image_url_https": "https://si0.twimg.com/profile_images/2284174872/7df3h38zabcvjylnyfe3_normal.png",
      "utc_offset": -28800,
      "id": 6253282,
      "profile_use_background_image": true,
      "listed_count": 10775,
      "profile_text_color": "333333",
      "lang": "en",
      "followers_count": 1212864,
      "protected": false,
      "notifications": null,
      "profile_background_image_url_https": "https://si0.twimg.com/images/themes/theme1/bg.png",
      "profile_background_color": "C0DEED",
      "verified": true,
      "geo_enabled": true,
      "time_zone": "Pacific Time (US & Canada)",
      "description": "The Real Twitter API. I tweet about API changes, service issues and happily answer questions about Twitter and our API. Don't get an answer? It's on my website.",
      "default_profile_image": false,
      "profile_background_image_url": "http://a0.twimg.com/images/themes/theme1/bg.png",
      "statuses_count": 3333,
      "friends_count": 31,
      "following": null,
      "show_all_inline_media": false,
      "screen_name": "twitterapi"
    },
    "in_reply_to_screen_name": null,
    "source": "<a href=\"http://sites.google.com/site/yorufukurou/\" rel=\"nofollow\">YoruFukurou</a>",
    "in_reply_to_status_id": null
  },
  {
    "coordinates": null,
    "favorited": false,
    "truncated": false,
    "created_at": "Sat Aug 25 17:26:51 +0000 2012",
    "id_str": "239413543487819778",
    "entities": {
      "urls": [
        {
          "expanded_url": "https://dev.twitter.com/issues/485",
          "url": "https://t.co/p5bOzH0k",
          "indices": [
            97,
            118
          ],
          "display_url": "dev.twitter.com/issues/485"
        }
      ],
      "hashtags": [

      ],
      "user_mentions": [

      ]
    },
    "in_reply_to_user_id_str": null,
    "contributors": null,
    "text": "We are working to resolve issues with application management &amp; logging in to the dev portal: https://t.co/p5bOzH0k ^TS",
    "retweet_count": 105,
    "in_reply_to_status_id_str": null,
    "id": 239413543487819778,
    "geo": null,
    "retweeted": false,
    "possibly_sensitive": false,
    "in_reply_to_user_id": null,
    "place": null,
    "user": {
      "profile_sidebar_fill_color": "DDEEF6",
      "profile_sidebar_border_color": "C0DEED",
      "profile_background_tile": false,
      "name": "Twitter API",
      "profile_image_url": "http://a0.twimg.com/profile_images/2284174872/7df3h38zabcvjylnyfe3_normal.png",
      "created_at": "Wed May 23 06:01:13 +0000 2007",
      "location": "San Francisco, CA",
      "follow_request_sent": false,
      "profile_link_color": "0084B4",
      "is_translator": false,
      "id_str": "6253282",
      "entities": {
        "url": {
          "urls": [
            {
              "expanded_url": null,
              "url": "http://dev.twitter.com",
              "indices": [
                0,
                22
              ]
            }
          ]
        },
        "description": {
          "urls": [

          ]
        }
      },
      "default_profile": true,
      "contributors_enabled": true,
      "favourites_count": 24,
      "url": "http://dev.twitter.com",
      "profile_image_url_https": "https://si0.twimg.com/profile_images/2284174872/7df3h38zabcvjylnyfe3_normal.png",
      "utc_offset": -28800,
      "id": 6253282,
      "profile_use_background_image": true,
      "listed_count": 10775,
      "profile_text_color": "333333",
      "lang": "en",
      "followers_count": 1212864,
      "protected": false,
      "notifications": null,
      "profile_background_image_url_https": "https://si0.twimg.com/images/themes/theme1/bg.png",
      "profile_background_color": "C0DEED",
      "verified": true,
      "geo_enabled": true,
      "time_zone": "Pacific Time (US & Canada)",
      "description": "The Real Twitter API. I tweet about API changes, service issues and happily answer questions about Twitter and our API. Don't get an answer? It's on my website.",
      "default_profile_image": false,
      "profile_background_image_url": "http://a0.twimg.com/images/themes/theme1/bg.png",
      "statuses_count": 3333,
      "friends_count": 31,
      "following": null,
      "show_all_inline_media": false,
      "screen_name": "twitterapi"
    },
    "in_reply_to_screen_name": null,
    "source": "<a href=\"http://sites.google.com/site/yorufukurou/\" rel=\"nofollow\">YoruFukurou</a>",
    "in_reply_to_status_id": null
  }
]`)

func BenchmarkParser(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var p Parser
		var r, n int

		for r < len(sample) {
			n, _, _ = p.Next(sample[r:])
			r += n
		}
	}
}

// step is a function which performs part of a test case.
type step func(*Parser, []byte) (int, bool, string)

// emit generates a function which invokes Parser.Next(), making sure the
// outcome matches what is expected.
func emit(offset int, want Event) step {
	return func(p *Parser, in []byte) (int, bool, string) {
		n, actual, err := p.Next(in)
		log := fmt.Sprintf(".Next(%#q) -> %d, %s, %#v", in, n, actual, err)

		if n != offset || actual != want {
			log += fmt.Sprintf(" (want %d, %s, <nil>)", offset, want)
			return n, false, log
		}

		return n, true, log
	}
}

// end generates a function which invokes Parser.End(), making sure the
// outcome matches what is expected.
func end(want Event) step {
	return func(p *Parser, in []byte) (int, bool, string) {
		actual, err := p.End()
		log := fmt.Sprintf(".End() -> %s, %#v", actual, err)

		if actual != want {
			log += fmt.Sprintf(" (want %s, <nil>)", want)
			return 0, false, log
		}

		return 0, true, log
	}
}

// parserTest describes a parser test case.
type parserTest struct {
	json  string
	steps []step
}

var parserTests = []parserTest{
	{
		`""`,
		[]step{
			emit(1, StringStart),
			emit(1, StringEnd),
			end(Done),
		},
	},
	{
		`"abc"`,
		[]step{
			emit(1, StringStart),
			emit(4, StringEnd),
			end(Done),
		},
	},
	{
		`"\u8bA0"`,
		[]step{
			emit(1, StringStart),
			emit(7, StringEnd),
			end(Done),
		},
	},
	{
		`"\u12345"`,
		[]step{
			emit(1, StringStart),
			emit(8, StringEnd),
			end(Done),
		},
	},
	{
		"\"\\b\\f\\n\\r\\t\\\\\"",
		[]step{
			emit(1, StringStart),
			emit(13, StringEnd),
			end(Done),
		},
	},
	{
		`0`,
		[]step{
			emit(1, NumberStart),
			end(NumberEnd),
		},
	},
	{
		`123`,
		[]step{
			emit(1, NumberStart),
			emit(2, Continue),
			end(NumberEnd),
		},
	},
	{
		`-10`,
		[]step{
			emit(1, NumberStart),
			emit(2, Continue),
			end(NumberEnd),
		},
	},
	{
		`0.10`,
		[]step{
			emit(1, NumberStart),
			emit(3, Continue),
			end(NumberEnd),
		},
	},
	{
		`12.03e+1`,
		[]step{
			emit(1, NumberStart),
			emit(7, Continue),
			end(NumberEnd),
		},
	},
	{
		`0.1e2`,
		[]step{
			emit(1, NumberStart),
			emit(4, Continue),
			end(NumberEnd),
		},
	},
	{
		`1e1`,
		[]step{
			emit(1, NumberStart),
			emit(2, Continue),
			end(NumberEnd),
		},
	},
	{
		`3.141569`,
		[]step{
			emit(1, NumberStart),
			emit(7, Continue),
			end(NumberEnd),
		},
	},
	{
		`10000000000000e-10`,
		[]step{
			emit(1, NumberStart),
			emit(17, Continue),
			end(NumberEnd),
		},
	},
	{
		`9223372036854775808`,
		[]step{
			emit(1, NumberStart),
			emit(18, Continue),
			end(NumberEnd),
		},
	},
	{
		`6E-06`,
		[]step{
			emit(1, NumberStart),
			emit(4, Continue),
			end(NumberEnd),
		},
	},
	{
		`1E-06`,
		[]step{
			emit(1, NumberStart),
			emit(4, Continue),
			end(NumberEnd),
		},
	},
	{
		`false`,
		[]step{
			emit(1, BoolStart),
			emit(4, BoolEnd),
			end(Done),
		},
	},
	{
		`true`,
		[]step{
			emit(1, BoolStart),
			emit(3, BoolEnd),
			end(Done),
		},
	},
	{
		`null`,
		[]step{
			emit(1, NullStart),
			emit(3, NullEnd),
			end(Done),
		},
	},
	{
		`"`,
		[]step{
			emit(1, StringStart),
			end(SyntaxError),
		},
	},
	{
		`"foo`,
		[]step{
			emit(1, StringStart),
			emit(3, Continue),
			end(SyntaxError),
		},
	},
	{
		`'single'`,
		[]step{
			emit(0, SyntaxError),
		},
	},
	{
		`"\u12g8"`,
		[]step{
			emit(1, StringStart),
			emit(4, SyntaxError),
		},
	},
	{
		`"\u"`,
		[]step{
			emit(1, StringStart),
			emit(2, SyntaxError),
		},
	},
	{
		`"you can\'t do this"`,
		[]step{
			emit(1, StringStart),
			emit(8, SyntaxError),
		},
	},
	{
		`-`,
		[]step{
			emit(1, NumberStart),
			end(SyntaxError),
		},
	},
	{
		`0.`,
		[]step{
			emit(1, NumberStart),
			emit(1, Continue),
			end(SyntaxError),
		},
	},
	{
		`123.456.789`,
		[]step{
			emit(1, NumberStart),
			emit(6, NumberEnd),
			emit(0, SyntaxError),
		},
	},
	{
		`10e`,
		[]step{
			emit(1, NumberStart),
			emit(2, Continue),
			end(SyntaxError),
		},
	},
	{
		`10e+`,
		[]step{
			emit(1, NumberStart),
			emit(3, Continue),
			end(SyntaxError),
		},
	},
	{
		`10e-`,
		[]step{
			emit(1, NumberStart),
			emit(3, Continue),
			end(SyntaxError),
		},
	},
	{
		`0e1x`,
		[]step{
			emit(1, NumberStart),
			emit(2, NumberEnd),
			emit(0, SyntaxError),
		},
	},
	{
		`0e+13.`,
		[]step{
			emit(1, NumberStart),
			emit(4, NumberEnd),
			emit(0, SyntaxError),
		},
	},
	{
		`0e+-0`,
		[]step{
			emit(1, NumberStart),
			emit(2, SyntaxError),
		},
	},
	{
		`tr`,
		[]step{
			emit(1, BoolStart),
			emit(1, Continue),
			end(SyntaxError),
		},
	},
	{
		`truE`,
		[]step{
			emit(1, BoolStart),
			emit(2, SyntaxError),
		},
	},
	{
		`fals`,
		[]step{
			emit(1, BoolStart),
			emit(3, Continue),
			end(SyntaxError),
		},
	},
	{
		`fALSE`,
		[]step{
			emit(1, BoolStart),
			emit(0, SyntaxError),
		},
	},
	{
		`n`,
		[]step{
			emit(1, NullStart),
			end(SyntaxError),
		},
	},
	{
		`NULL`,
		[]step{
			emit(0, SyntaxError),
		},
	},
	{
		`{}`,
		[]step{
			emit(1, ObjectStart),
			emit(1, ObjectEnd),
			end(Done),
		},
	},
	{
		`{"foo":"bar"}`,
		[]step{
			emit(1, ObjectStart),
			emit(1, KeyStart),
			emit(4, KeyEnd),
			emit(2, StringStart),
			emit(4, StringEnd),
			emit(1, ObjectEnd),
			end(Done),
		},
	},
	{
		`{"a":1`,
		[]step{
			emit(1, ObjectStart),
			emit(1, KeyStart),
			emit(2, KeyEnd),
			emit(2, NumberStart),
			end(SyntaxError),
		},
	},
	{
		`{"1":1,"2":2}`,
		[]step{
			emit(1, ObjectStart),
			emit(1, KeyStart),
			emit(2, KeyEnd),
			emit(2, NumberStart),
			emit(0, NumberEnd),
			emit(2, KeyStart),
			emit(2, KeyEnd),
			emit(2, NumberStart),
			emit(0, NumberEnd),
			emit(1, ObjectEnd),
			end(Done),
		},
	},
	{
		`{"\u1234\t\n\b\u8BbF\"":0}`,
		[]step{
			emit(1, ObjectStart),
			emit(1, KeyStart),
			emit(21, KeyEnd),
			emit(2, NumberStart),
			emit(0, NumberEnd),
			emit(1, ObjectEnd),
			end(Done),
		},
	},
	{
		`{"{":"}"}`,
		[]step{
			emit(1, ObjectStart),
			emit(1, KeyStart),
			emit(2, KeyEnd),
			emit(2, StringStart),
			emit(2, StringEnd),
			emit(1, ObjectEnd),
			end(Done),
		},
	},
	{
		`{"foo":{"bar":{"baz":{}}}}`,
		[]step{
			emit(1, ObjectStart),
			emit(1, KeyStart),
			emit(4, KeyEnd),
			emit(2, ObjectStart),
			emit(1, KeyStart),
			emit(4, KeyEnd),
			emit(2, ObjectStart),
			emit(1, KeyStart),
			emit(4, KeyEnd),
			emit(2, ObjectStart),
			emit(1, ObjectEnd),
			emit(1, ObjectEnd),
			emit(1, ObjectEnd),
			emit(1, ObjectEnd),
			end(Done),
		},
	},
	{
		`{"true":true,"false":false,"null":null}`,
		[]step{
			emit(1, ObjectStart),
			emit(1, KeyStart),
			emit(5, KeyEnd),
			emit(2, BoolStart),
			emit(3, BoolEnd),
			emit(2, KeyStart),
			emit(6, KeyEnd),
			emit(2, BoolStart),
			emit(4, BoolEnd),
			emit(2, KeyStart),
			emit(5, KeyEnd),
			emit(2, NullStart),
			emit(3, NullEnd),
			emit(1, ObjectEnd),
			end(Done),
		},
	},
	{
		`{0:1}`,
		[]step{
			emit(1, ObjectStart),
			emit(0, SyntaxError),
		},
	},
	{
		`{"foo":"bar"`,
		[]step{
			emit(1, ObjectStart),
			emit(1, KeyStart),
			emit(4, KeyEnd),
			emit(2, StringStart),
			emit(4, StringEnd),
			end(SyntaxError),
		},
	},
	{
		`{{`,
		[]step{
			emit(1, ObjectStart),
			emit(0, SyntaxError),
		},
	},
	{
		`{"a":1,}`,
		[]step{
			emit(1, ObjectStart),
			emit(1, KeyStart),
			emit(2, KeyEnd),
			emit(2, NumberStart),
			emit(0, NumberEnd),
			emit(1, SyntaxError),
		},
	},
	{
		`{"a":1,,`,
		[]step{
			emit(1, ObjectStart),
			emit(1, KeyStart),
			emit(2, KeyEnd),
			emit(2, NumberStart),
			emit(0, NumberEnd),
			emit(1, SyntaxError),
		},
	},
	{
		`{"a"}`,
		[]step{
			emit(1, ObjectStart),
			emit(1, KeyStart),
			emit(2, KeyEnd),
			emit(0, SyntaxError),
		},
	},
	{
		`{"a":"1}`,
		[]step{
			emit(1, ObjectStart),
			emit(1, KeyStart),
			emit(2, KeyEnd),
			emit(2, StringStart),
			emit(2, Continue),
			end(SyntaxError),
		},
	},
	{
		`{"a":1"b":2}`,
		[]step{
			emit(1, ObjectStart),
			emit(1, KeyStart),
			emit(2, KeyEnd),
			emit(2, NumberStart),
			emit(0, NumberEnd),
			emit(0, SyntaxError),
		},
	},

	{
		`[]`,
		[]step{
			emit(1, ArrayStart),
			emit(1, ArrayEnd),
			end(Done),
		},
	},
	{
		`[1]`,
		[]step{
			emit(1, ArrayStart),
			emit(1, NumberStart),
			emit(0, NumberEnd),
			emit(1, ArrayEnd),
			end(Done),
		},
	},
	{
		`[1,2,3]`,
		[]step{
			emit(1, ArrayStart),
			emit(1, NumberStart),
			emit(0, NumberEnd),
			emit(2, NumberStart),
			emit(0, NumberEnd),
			emit(2, NumberStart),
			emit(0, NumberEnd),
			emit(1, ArrayEnd),
			end(Done),
		},
	},
	{
		`["dude","what"]`,
		[]step{
			emit(1, ArrayStart),
			emit(1, StringStart),
			emit(5, StringEnd),
			emit(2, StringStart),
			emit(5, StringEnd),
			emit(1, ArrayEnd),
			end(Done),
		},
	},
	{
		`[[[],[]],[]]`,
		[]step{
			emit(1, ArrayStart),
			emit(1, ArrayStart),
			emit(1, ArrayStart),
			emit(1, ArrayEnd),
			emit(2, ArrayStart),
			emit(1, ArrayEnd),
			emit(1, ArrayEnd),
			emit(2, ArrayStart),
			emit(1, ArrayEnd),
			emit(1, ArrayEnd),
			end(Done),
		},
	},
	{
		`[10.3]`,
		[]step{
			emit(1, ArrayStart),
			emit(1, NumberStart),
			emit(3, NumberEnd),
			emit(1, ArrayEnd),
			end(Done),
		},
	},
	{
		`["][\"]"]`,
		[]step{
			emit(1, ArrayStart),
			emit(1, StringStart),
			emit(6, StringEnd),
			emit(1, ArrayEnd),
			end(Done),
		},
	},
	{
		`[`,
		[]step{
			emit(1, ArrayStart),
			end(SyntaxError),
		},
	},
	{
		`[,]`,
		[]step{
			emit(1, ArrayStart),
			emit(0, SyntaxError),
		},
	},
	{
		`[10,]`,
		[]step{
			emit(1, ArrayStart),
			emit(1, NumberStart),
			emit(1, NumberEnd),
			emit(1, SyntaxError),
		},
	},
	{
		`[}`,
		[]step{
			emit(1, ArrayStart),
			emit(0, SyntaxError),
		},
	},
	{
		` 17`,
		[]step{
			emit(2, NumberStart),
			emit(1, Continue),
			end(NumberEnd),
		},
	},
	{`38 `,
		[]step{
			emit(1, NumberStart),
			emit(1, NumberEnd),
			emit(1, Continue),
			end(Done),
		},
	},
	{
		`  " what ? "  `,
		[]step{
			emit(3, StringStart),
			emit(9, StringEnd),
			emit(2, Continue),
			end(Done),
		},
	},
	{
		"\nnull",
		[]step{
			emit(2, NullStart),
			emit(3, NullEnd),
			end(Done),
		},
	},
	{
		"\n\r\t true \r\n\t",
		[]step{
			emit(5, BoolStart),
			emit(3, BoolEnd),
			emit(4, Continue),
			end(Done),
		},
	},
	{
		" { \"foo\": \t\"bar\" } ",
		[]step{
			emit(2, ObjectStart),
			emit(2, KeyStart),
			emit(4, KeyEnd),
			emit(4, StringStart),
			emit(4, StringEnd),
			emit(2, ObjectEnd),
			emit(1, Continue),
			end(Done),
		},
	},
	{
		"\t[ 1 , 2\r, 3\n]",
		[]step{
			emit(2, ArrayStart),
			emit(2, NumberStart),
			emit(0, NumberEnd),
			emit(4, NumberStart),
			emit(0, NumberEnd),
			emit(4, NumberStart),
			emit(0, NumberEnd),
			emit(2, ArrayEnd),
			end(Done),
		},
	},
	{
		" \n{ \t\"foo\" : [ \"bar\", null ], \"what\\n\"\t: 10.3e1 } ",
		[]step{
			emit(3, ObjectStart),
			emit(3, KeyStart),
			emit(4, KeyEnd),
			emit(4, ArrayStart),
			emit(2, StringStart),
			emit(4, StringEnd),
			emit(3, NullStart),
			emit(3, NullEnd),
			emit(2, ArrayEnd),
			emit(3, KeyStart),
			emit(7, KeyEnd),
			emit(4, NumberStart),
			emit(5, NumberEnd),
			emit(2, ObjectEnd),
			emit(1, Continue),
			end(Done),
		},
	},
}

func TestParser(t *testing.T) {
	for _, test := range parserTests {
		b := []byte(test.json)
		p := &Parser{}

		// buffer each step's log output, and print them all if
		// the test case fails
		output := "\np := Parser{}"

		for i := 0; i < len(test.steps); i++ {
			n, ok, log := test.steps[i](p, b)
			output += "\np" + log

			if !ok {
				t.Errorf(output)
				break
			}

			b = b[n:]
		}
	}
}
