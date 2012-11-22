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

var Post = Backbone.Model.extend({
	initialize : function() {
		this.author = ""
		this.link = ""
		this.text = ""
		this.timestamp = ""
		this.comments = new Comments();
	},
	comparator : function(post) {
		return post.get('timestamp')
	}
});

var Posts = Backbone.Collection.extend({
	model: Post
});

var CommentView = Backbone.View.extend({
	className: "commentView",
	tagName: "li",
	render : function() {
		this.$el.html(post_template({model : this.model}));
		return this;
	}
});

var PostView = Backbone.View.extend({
	className : "post",
	tagName : "li",
	render : function() {
		//this.$el.html("<p>" + this.model.get("text") + '<p>');
		console.log("model:");
		console.log(this.model);
		console.log(this.model.get('text'));
		this.$el.html(post_template({model : this.model}));
		return this;
	}

});

var Posts = new Posts();

var AppView = Backbone.View.extend({
	el : "#messages",
	addOne : function(post) {
		var postView = new PostView({model: post});
		console.log('post:');
		console.log(post);
		console.log(postView);
		console.log(this.$("#message-list"));
		$("#message_list").append(postView.render().el);
	},
	initialize : function() {
		Posts.bind('add', this.addOne, this);
	}

});

var App = new AppView();

var init = function() {
	console.log("init");
	for (var i = 0; i < 10; i++) {
		var post = new Post({author: "mislav", text: "Remember remember the fifth of november", timestamp: i});
		Posts.add(post);
	}
};
$(init);

var post_template = _.template(
	'<p><%model.get("text")%></p>' +
	'<p><%model.get("timestamp")%></p>'+
	'<p><%-"aaa"%></p>'
	);
