// Copyright © Rob Burke inchworks.com, 2020.

// This file is part of PicInch.
//
// PicInch is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// PicInch is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with PicInch.  If not, see <https://www.gnu.org/licenses/>.

package models

// Database models PicInch.

import (
	"errors"
	"html/template"
	"strconv"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// Database field names are the same as structure names, with lower case first letter.

const (
	// slide formats
	SlideBlank = iota
	SlideT
	SlideI
	SlideTI
	SlideC
	SlideTC
	SlideIC
	SlideTIC

	// slideshow visibility
	SlideshowTopic   = -1
	SlideshowPrivate = 0
	SlideshowClub    = 1
	SlideshowPublic  = 2

	// user status
	UserSuspended = 0
	UserKnown     = 1
	UserActive    = 2
	UserCurator   = 3
	UserAdmin     = 4
)

var (
	ErrNoRecord           = errors.New("models: no matching record found")
	ErrInvalidCredentials = errors.New("models: invalid credentials")
	ErrDuplicateEmail     = errors.New("models: duplicate email")
)

var VisibleOpts = []string{"none", "club", "public"}

type Gallery struct {
	Id int64

	// parameters
	Organiser  string
	NMaxSlides int `db:"n_max_slides"`
	NShowcased int `db:"n_showcased"`

	// announcements
	NoticePublic string // appears on home page
	NoticeUsers  string // appears on contributor's page
}

type Slide struct {
	Id        int64
	Slideshow int64
	Format    int
	ShowOrder int `db:"show_order"`
	Created   time.Time
	Revised   time.Time
	Title     string
	Caption   string
	Image     string
}

type Slideshow struct {
	Id           int64
	Gallery      int64
	GalleryOrder int `db:"gallery_order"`
	Visible      int
	User         int64
	Shared       int64 // link for external access
	Topic        int64 // 0 for a normal slideshow
	Created      time.Time
	Revised      time.Time
	Title        string
	Caption      string
	Format       string
	Image        string
}

type Topic struct {
	Id           int64
	Gallery      int64
	GalleryOrder int `db:"gallery_order"`
	Visible      int
	Shared       int64 // link for external access
	Created      time.Time
	Revised      time.Time
	Title        string
	Caption      string
	Format       string
	Image        string
}

type User struct {
	Id      int64
	Gallery int64
	// Shared       int64  // ## link for external access
	Username string
	Name     string

	// user management
	Status   int
	Password []byte
	Created  time.Time
}

// Join results

type SlideshowUser struct {
	Id    int64
	Title string
	Image string
	Name  string // user's display name
}

type TopicSlide struct {
	Format  int
	Title   string
	Caption string
	Image   string
	Name    string
}

// Fields with newlines replaced by breaks, and HTML formatting allowed.
// ## If source is untrusted, could return a slice of lines and use range to add breaks in template.

func (p *Slide) TitleBr() template.HTML {
	return Nl2br(p.Title)
}

func (p *Slide) CaptionBr() template.HTML {
	return Nl2br(p.Caption)
}

func Nl2br(str string) template.HTML {
	return template.HTML(strings.Replace(str, "\n", "<br>", -1))
}

// Code to string conversions

func (s *Slideshow) VisibleStr() string {

	return VisibleOpts[s.Visible]
}

func (s *Topic) VisibleStr() string {

	return VisibleOpts[s.Visible]
}

// Formats

func (t *Topic) ParseFormat() (fmt string, max int) {

	var err error
	ss := strings.Split(t.Format, ".")

	switch len(ss) {
	case 0:
		return "T", 8

	case 1:
		fmt = ss[0]
		max = 8

	default:
		fmt = ss[0]
		max, err = strconv.Atoi(ss[1])
		if err != nil {
			max = 8
		} // default
	}

	return
}

// Password management

func (u *User) Authenticate(pwd string) error {

	// must be an active user
	if u.Status < UserActive {
		return ErrInvalidCredentials
	}

	// check password
	err := bcrypt.CompareHashAndPassword(u.Password, []byte(pwd))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return ErrInvalidCredentials
		} else {
			return err
		}
	}
	return nil
}

func (u *User) SetPassword(pwd string) error {

	hashed, err := bcrypt.GenerateFromPassword([]byte(pwd), 12)
	if err != nil {
		return err
	} else {
		u.Password = hashed
	}
	return nil
}
