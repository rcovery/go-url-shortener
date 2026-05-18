package shorturl

import "net/url"

type Link url.URL

func NewLink(rawURL string) (*Link, error) {
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return nil, err
	}

	// This pointer (parsedURL) was associated with url.URL
	// Now, I'm associating it with another pointer with type Link
	// So, at the end, I'm just converting pointer types
	return (*Link)(parsedURL), nil
}

func (l *Link) String() string {
	// The same from above:
	return (*url.URL)(l).String()
}

func (l *Link) Equals(anotherLink *Link) bool {
	return l.String() == anotherLink.String()
}

/*
Alternatively, I could do this:

func (l Link) String() string {
	return (*url.URL)(&l).String()
}

It's just a design decision
*/
