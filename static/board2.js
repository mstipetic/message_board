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
	idAttribute: 'id',
	url: function() {
		return '/post' + (this.id ? '/' + this.id : '');
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
//	url : '/post'
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
		comment.url = this.model.url() + '/comments';
		console.log(comment.url);
		var that = this;
		$.when(comment.save()).then(function() {
			that.$('textarea').val('');
			that.model.get('comments').add(comment);
		});
	}
});

var NewPostView = Backbone.View.extend({
	el : "#new_message",
	render: function() {
		console.log(this.$el);
		console.log('render');
		this.$el.append($(new_message_template()));
	},
	events: {
		'click #new_message_btn' : 'new_post',
	},
	new_post : function() {
		console.log('new post');
		var title = this.$('#new_message_title').val();
		var url = this.$('#new_message_url').val();
		var text = this.$('#new_message_text').val();
		this.$('#new_message_title').val('');
		this.$('#new_message_url').val('');
		this.$('#new_message_text').val('');
		var post = new Post({title: title, url: url, text: text});
		$.when(post.save()).then(function() {
			posts.add(post);
			//console.log('hej');
			//console.log(post.id);
		});
	}
});

var App = Backbone.View.extend({
	el : "#messages",
	initialize : function() {
		posts.bind('add', this.addPost, this);
		this.getPosts();
	},
	addPost : function(post) {
		var postView = new PostView({model : post});
		this.$el.append(postView.render().$el);
	},
	getPosts : function(page) {
		page = page || 0;
		$.when($.get('/post', {page : page}))
			.then(function(post_list) { 
				console.log('u callbacku');
				_.each(post_list, function(post) { posts.add(post) });
				console.log(posts);
			});
	}
});



$(function() {
	var app = new App();
	var newPostView= new NewPostView();
	//console.log(newPostView);
	newPostView.render();
	for (var i = 0; i < 0; i++) {
		$('#new_message_title').val('Message title');
		$('#new_message_url').val('http://www.google.com');
		$('#new_message_text').val('svasta nesta kul');
		$('#new_message_btn').click();
	}
});


var post_template = _.template($('#post_template').html());
var comment_template = _.template($('#comment_template').html());
var new_message_template = _.template($('#new_message_template').html());
