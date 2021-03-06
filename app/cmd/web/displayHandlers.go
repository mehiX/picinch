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

package main

// Requests for gallery display pages

import (
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"

	"inchworks.com/gallery/pkg/usage"
)

// About page for website

func (app *Application) about(w http.ResponseWriter, r *http.Request) {

	data := DataCommon{}
	data.addDefaultData(app, r)

	app.render(w, r, "about.page.tmpl", &data)
}

// Contributor (for other users to see)

func (app *Application) contributor(w http.ResponseWriter, r *http.Request) {

	ps := httprouter.ParamsFromContext(r.Context())
	userId, _ := strconv.ParseInt(ps.ByName("nUser"), 10, 64)

	// template and data for contributor
	template, data := app.galleryState.DisplayContributor(userId)
	if data == nil {
		// ## better to show "unknown contributor" nicely
		app.clientError(w, http.StatusBadRequest)
		return
	}

	// display page
	app.render(w, r, template, data)
}

// Contributors (for other users to see)

func (app *Application) contributors(w http.ResponseWriter, r *http.Request) {

	// template and contributors
	template, data := app.galleryState.DisplayContributors()

	// display page
	app.render(w, r, template, data)
}

// Highlighted image, to be embedded in parent website

func (app *Application) highlight(w http.ResponseWriter, r *http.Request) {

	ps := httprouter.ParamsFromContext(r.Context())
	prefix := ps.ByName("prefix")
	n, _ := strconv.Atoi(ps.ByName("nImage"))

	// get highlighted image
	image := app.galleryState.Highlighted(prefix, n)

	// return image
	if image != "" {
		http.ServeFile(w, r, image)
	} else {
		app.notFound(w)
	}
}

// Highlights, to be embedded in parent website

func (app *Application) highlights(w http.ResponseWriter, r *http.Request) {

	ps := httprouter.ParamsFromContext(r.Context())
	nImages, _ := strconv.Atoi(ps.ByName("nImages"))

	data := app.galleryState.DisplayEmbedded(nImages)

	app.render(w, r, "highlights.page.tmpl", data)
}

// Home page for website

func (app *Application) home(w http.ResponseWriter, r *http.Request) {

	template, data := app.galleryState.DisplayHome(app.isAuthenticated(r))
	if data == nil {
		app.clientError(w, http.StatusInternalServerError)
		return
	}

	app.render(w, r, template, data)
}

// Logout user

func (app *Application) logout(w http.ResponseWriter, r *http.Request) {

	// remove user ID from the session data
	app.session.Remove(r, "authenticatedUserID")

	// flash message to confirm logged out
	app.session.Put(r, "flash", "You are logged out")

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// Slideshow

func (app *Application) slideshow(w http.ResponseWriter, r *http.Request) {

	ps := httprouter.ParamsFromContext(r.Context())

	showId, _ := strconv.ParseInt(ps.ByName("nShow"), 10, 64)

	// allow access to show?
	// ## reads show, and DisplaySlideshow will read it again
	if !app.allowViewShow(r, showId) {
		app.clientError(w, http.StatusUnauthorized)
		return
	}

	// template and data for slides
	template, data := app.galleryState.DisplaySlideshow(showId, r.Referer())
	if data == nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	// display page
	app.render(w, r, template, data)
}

// Slideshows for user

func (app *Application) slideshowsOwn(w http.ResponseWriter, r *http.Request) {

	// user
	userId := app.authenticatedUser(r)

	data := app.galleryState.ForMyGallery(userId)

	app.render(w, r, "my-gallery.page.tmpl", data)
}

func (app *Application) slideshowsUser(w http.ResponseWriter, r *http.Request) {

	ps := httprouter.ParamsFromContext(r.Context())
	userId, _ := strconv.ParseInt(ps.ByName("nUser"), 10, 64)

	// allow access?
	if !app.allowAccessUser(r, userId) {
		app.clientError(w, http.StatusUnauthorized)
		return
	}

	data := app.galleryState.ForMyGallery(userId)

	app.render(w, r, "my-gallery.page.tmpl", data)
}

// All slideshows for topic

func (app *Application) topic(w http.ResponseWriter, r *http.Request) {

	ps := httprouter.ParamsFromContext(r.Context())

	id, _ := strconv.ParseInt(ps.ByName("nShow"), 10, 64)
	seq, _ := strconv.ParseInt(ps.ByName("seq"), 10, 32)

	// allow access to topic?
	// ## reads topic, and DisplayTopic will read it again
	if !app.allowViewTopic(r, id) {
		app.clientError(w, http.StatusUnauthorized)
		return
	}

	// template and data for slides
	template, data := app.galleryState.DisplayTopic(id, int(seq), "/")
	if data == nil {
		app.session.Put(r, "flash", "No contributions to this topic yet.")
		http.Redirect(w, r, r.Referer(), http.StatusSeeOther)
		return
	}

	// display page
	app.render(w, r, template, data)
}

// Topic slides for user

func (app *Application) topicUser(w http.ResponseWriter, r *http.Request) {

	ps := httprouter.ParamsFromContext(r.Context())

	showId, _ := strconv.ParseInt(ps.ByName("nShow"), 10, 64)
	userId, _ := strconv.ParseInt(ps.ByName("nUser"), 10, 64)

	// allow access?
	if !app.allowAccessUser(r, userId) {
		app.clientError(w, http.StatusUnauthorized)
		return
	}

	// template and data for slides
	template, data := app.galleryState.DisplayTopicUser(showId, userId, r.Referer())
	if data == nil {
		app.session.Put(r, "flash", "No slides to this topic yet.")
		http.Redirect(w, r, r.Referer(), http.StatusSeeOther)
		return
	}

	// display page
	app.render(w, r, template, data)
}

// Users slideshows for topic

func (app *Application) topicContributors(w http.ResponseWriter, r *http.Request) {

	ps := httprouter.ParamsFromContext(r.Context())

	topicId, _ := strconv.ParseInt(ps.ByName("nTopic"), 10, 64)

	// template and data for slides
	template, data := app.galleryState.DisplayTopicContributors(topicId)

	// display page
	app.render(w, r, template, data)
}

// Topics

func (app *Application) topics(w http.ResponseWriter, r *http.Request) {

	// allow access?
	if !app.isCurator(r) {
		app.clientError(w, http.StatusUnauthorized)
		return
	}

	data := app.galleryState.ForTopics()

	app.render(w, r, "topics.page.tmpl", data)
}

// Usage statistics

func (app *Application) usageDays(w http.ResponseWriter, r *http.Request) {

	// allow access?
	if !app.isAdmin(r) {
		app.clientError(w, http.StatusUnauthorized)
		return
	}

	data := app.galleryState.ForUsage(usage.Day)

	app.render(w, r, "usage.page.tmpl", data)
}

func (app *Application) usageMonths(w http.ResponseWriter, r *http.Request) {

	// allow access?
	if !app.isAdmin(r) {
		app.clientError(w, http.StatusUnauthorized)
		return
	}

	data := app.galleryState.ForUsage(usage.Month)

	app.render(w, r, "usage.page.tmpl", data)
}

// For curator

func (app *Application) usersCurator(w http.ResponseWriter, r *http.Request) {

	// allow access?
	if !app.isCurator(r) {
		app.clientError(w, http.StatusUnauthorized)
		return
	}

	data := app.galleryState.ForUsers()

	app.render(w, r, "users-curator.page.tmpl", data)
}
