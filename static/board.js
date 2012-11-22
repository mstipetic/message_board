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
		console.log("model:");
		console.log(this.model);
		console.log(this.model.get('text'));
		this.$el.html(post_template({model : this.model}));
		return this;
	},
	test : function() {
		console.log(this);
	}

});

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
	initialize : function() {
		console.log("comment view:");
		console.log(this.model)
		this.model.comments.bind('add', this.addOne, this);
	},
	addOne : function(comment) {
		console.log(this.$(".commentHolder"))
	}
});

var AppView = Backbone.View.extend({
	el : "#messages",
	addOne : function(post) {
		var postView = new PostView({model: post});
		console.log(post);
		var commentView = new CommentView({model: post});
		$("#message_list").append(postView.render().el);
	},
	initialize : function() {
		posts.bind('add', this.addOne, this);
	}

});

var App = new AppView();

var init = function() {
	console.log("init");
	for (var i = 0; i < 10; i++) {
		var post = new Post({author: "mislav", text: "Remember remember the fifth of november", timestamp: i});
		posts.add(post);
		var comment = new Comment{{author : 'mislav', text: 'com'});
		
	}
};
$(init);

var post_template = _.template(
	'<div class="postAuthor"><%print(model.get("author"))%></div>' +
	'<div class="postText"><%print(model.get("text"))%></div>' +
	'<div class="postTime"><%print(model.get("timestamp"))%></div>' +
	'<hr />' +
	'<div class="commentHolder"></div>'
	);
