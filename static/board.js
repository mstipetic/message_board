//************************************************************************
var Comment = Backbone.Model.extend({
	initialize : function() {
		this.author = ""
		this.text = ""
		this.timestamp = ""
	}
});

var Comments = Backbone.Collection.extend({
	model: Comment
});

var CommentView = Backbone.View.extend({
	tagName : "li",
	className : "comment",
	render: function() {
		this.$el.html('<p>comment</p>');
		return this;
	}
});

var CommentsApp = Backbone.View.extend({
	initialize: function() {
		this.model.comments.bind('add', this.addOne, this);
	},
	addOne : function(comment) {
		console.log('comments app add one:');
		//console.log(this)
		//console.log($(".commentHolder"))
		var commentView = new CommentView({model: comment})
		this.$el.append(commentView.render().$el);
	},
});
// ***********************************************************************
var Post = Backbone.Model.extend({
	initialize : function() {
	},
	comments : new Comments(),
	comparator : function(post) {
		return post.get('timestamp')
	}
});

var Posts = Backbone.Collection.extend({
	model: Post
});

var posts = new Posts();

var PostView = Backbone.View.extend({
	className : "post",
	tagName : "li",
	events: {
		"click"	: "test"

	},
	render : function() {
		//this.$el.html("<p>" + this.model.get("text") + '<p>');
		console.log("postview render");
		//console.log("model:");
		//console.log(this.model);
		//console.log(this.model.get('text'));
		this.$el.html(post_template({model : this.model}));
		return this;
	},
	test : function() {
		console.log(this);
	}

});


var PostsApp = Backbone.View.extend({
	el : "#messages",
	addOne : function(post) {
		var postView = new PostView({model: post});
		console.log("posts app");
		console.log(post);
		postView.render();
		var commentApp = new CommentsApp({model: post, el : $(".commentList", postView.$el)[0]});
		console.log($(".commentList", postView.$el));
		$("#message_list").append(postView.el);
	},
	initialize : function() {
		posts.bind('add', this.addOne, this);
	}

});

var App = new PostsApp();

var init = function() {
	console.log("init");
	for (var i = 0; i < 5; i++) {
		var post = new Post({author: "mislav", text: "Remember remember the fifth of november", timestamp: i});
		posts.add(post);
		var comment = new Comment({author : 'mislav', text: 'com'});
		post.comments.add(comment);
		
	}
};
$(init);

var post_template = _.template(
	'<div class="postAuthor"><%print(model.get("author"))%></div>' +
	'<div class="postText"><%print(model.get("text"))%></div>' +
	'<div class="postTime"><%print(model.get("timestamp"))%></div>' +
	'<hr />' +
	'<div class="commentHolder"><ul class="commentList"></ul></div>'
	);
