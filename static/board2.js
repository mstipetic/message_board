var Comment = Backbone.RelationalModel.extend({
	url : function() {
		return this.get('post').url() + '/comments';
	}
});

var CommentView = Backbone.View.extend({
	tagName : 'li',
	className : 'comment',
	render : function() {
		this.$el.append($(comment_template({model : this.model})));
		return this;
	}
});

var Post = Backbone.RelationalModel.extend({
	initialize : function() {
	},
	relations : [{
		type : 'HasMany',
		key : 'comments',
		relatedModel: 'Comment',
		reverseRelation: {
			key : 'post'
		}
	}],
});

var PostCollection = Backbone.Collection.extend({
	model : Post,
	url : '/post'
});

var posts = new PostCollection();

var PostView = Backbone.View.extend({
	tagName : 'li',
	className: 'post',
	initialize: function() {
		this.model.bind('add:comments', this.addComment, this)
	},
	addComment : function(comment, attr) {
		var commentView = new CommentView({model : comment});
		console.log(comment.get('post'));
		this.$('.commentList').append(commentView.render().$el);
	},
	render : function() {
		this.$el.append($(post_template({model : this.model})));
		return this;
	},
	events: {
		'click input' : 'submitComment'
	},
	submitComment: function() {
		var comment = new Comment({text : this.$('textarea').val()});
		this.$('textarea').val('');
		this.model.get('comments').add(comment);
	}
});

var App = Backbone.View.extend({
	el : "#messages",
	initialize : function() {
		posts.bind('add', this.addPost, this);
	},
	addPost : function(post) {
		var postView = new PostView({model : post});
		this.$el.append(postView.render().$el);
	}
});



$(function() {
	var app = new App();
	for (var i = 0; i < 5; i++) {
		var post = new Post({text: "Remember remember the fifth of november", timestamp: i});
		posts.add(post);
		var comment= new Comment({author : 'mislav', text: 'com'});
		post.get('comments').add(comment);
		//console.log(post.get('comments'));
	}
});


var post_template = _.template($('#post_template').html());
var comment_template = _.template($('#comment_template').html());
