package bot

import (
	"github.com/Mrs4s/MiraiGo/client"
	"github.com/Mrs4s/MiraiGo/message"
)

// The following functions are designed to make the framework compatible with the old plug-in, and the newly developed plug-in should no longer use this method

// Deprecated: This function is designed to make the framework compatible with the old plug-in, and the newly developed plug-in should no longer use this method
func (bot *Bot) OnPrivateMessage(f func(qqClient *client.QQClient, event *message.PrivateMessage)) {
	bot.PrivateMessageEvent.Subscribe(f)
}

// Deprecated: This function is designed to make the framework compatible with the old plug-in, and the newly developed plug-in should no longer use this method
func (bot *Bot) OnPrivateMessageF(filter func(*message.PrivateMessage) bool, f func(*client.QQClient, *message.PrivateMessage)) {
	bot.OnPrivateMessage(func(client *client.QQClient, msg *message.PrivateMessage) {
		if filter(msg) {
			f(client, msg)
		}
	})
}

// Deprecated: This function is designed to make the framework compatible with the old plug-in, and the newly developed plug-in should no longer use this method
func (bot *Bot) OnTempMessage(f func(qqClient *client.QQClient, event *client.TempMessageEvent)) {
	bot.TempMessageEvent.Subscribe(f)
}

// Deprecated: This function is designed to make the framework compatible with the old plug-in, and the newly developed plug-in should no longer use this method
func (bot *Bot) OnGroupMessage(f func(qqClient *client.QQClient, event *message.GroupMessage)) {
	bot.GroupMessageEvent.Subscribe(f)
}

// Deprecated: This function is designed to make the framework compatible with the old plug-in, and the newly developed plug-in should no longer use this method
func (bot *Bot) OnSelfPrivateMessage(f func(qqClient *client.QQClient, event *message.PrivateMessage)) {
	bot.SelfPrivateMessageEvent.Subscribe(f)
}

// Deprecated: This function is designed to make the framework compatible with the old plug-in, and the newly developed plug-in should no longer use this method
func (bot *Bot) OnSelfGroupMessage(f func(qqClient *client.QQClient, event *message.GroupMessage)) {
	bot.SelfGroupMessageEvent.Subscribe(f)
}

// Deprecated: This function is designed to make the framework compatible with the old plug-in, and the newly developed plug-in should no longer use this method
func (bot *Bot) OnGroupMuted(f func(qqClient *client.QQClient, event *client.GroupMuteEvent)) {
	bot.GroupMuteEvent.Subscribe(f)
}

// Deprecated: This function is designed to make the framework compatible with the old plug-in, and the newly developed plug-in should no longer use this method
func (bot *Bot) OnGroupMessageRecalled(f func(qqClient *client.QQClient, event *client.GroupMessageRecalledEvent)) {
	bot.GroupMessageRecalledEvent.Subscribe(f)
}

// Deprecated: This function is designed to make the framework compatible with the old plug-in, and the newly developed plug-in should no longer use this method
func (bot *Bot) OnFriendMessageRecalled(f func(qqClient *client.QQClient, event *client.FriendMessageRecalledEvent)) {
	bot.FriendMessageRecalledEvent.Subscribe(f)
}

// Deprecated: This function is designed to make the framework compatible with the old plug-in, and the newly developed plug-in should no longer use this method
func (bot *Bot) OnGroupJoin(f func(qqClient *client.QQClient, event *client.GroupInfo)) {
	bot.GroupJoinEvent.Subscribe(f)
}

// Deprecated: This function is designed to make the framework compatible with the old plug-in, and the newly developed plug-in should no longer use this method
func (bot *Bot) OnGroupLeave(f func(qqClient *client.QQClient, event *client.GroupLeaveEvent)) {
	bot.GroupLeaveEvent.Subscribe(f)
}

// Deprecated: This function is designed to make the framework compatible with the old plug-in, and the newly developed plug-in should no longer use this method
func (bot *Bot) OnGroupMemberJoin(f func(qqClient *client.QQClient, event *client.MemberJoinGroupEvent)) {
	bot.GroupMemberJoinEvent.Subscribe(f)
}

// Deprecated: This function is designed to make the framework compatible with the old plug-in, and the newly developed plug-in should no longer use this method
func (bot *Bot) OnGroupMemberLeave(f func(qqClient *client.QQClient, event *client.MemberLeaveGroupEvent)) {
	bot.GroupMemberLeaveEvent.Subscribe(f)
}

// Deprecated: This function is designed to make the framework compatible with the old plug-in, and the newly developed plug-in should no longer use this method
func (bot *Bot) OnMemberCardUpdated(f func(qqClient *client.QQClient, event *client.MemberCardUpdatedEvent)) {
	bot.MemberCardUpdatedEvent.Subscribe(f)
}

// Deprecated: This function is designed to make the framework compatible with the old plug-in, and the newly developed plug-in should no longer use this method
func (bot *Bot) OnGroupNameUpdated(f func(qqClient *client.QQClient, event *client.GroupNameUpdatedEvent)) {
	bot.GroupNameUpdatedEvent.Subscribe(f)
}

// Deprecated: This function is designed to make the framework compatible with the old plug-in, and the newly developed plug-in should no longer use this method
func (bot *Bot) OnGroupMemberPermissionChanged(f func(qqClient *client.QQClient, event *client.MemberPermissionChangedEvent)) {
	bot.GroupMemberPermissionChangedEvent.Subscribe(f)
}

// Deprecated: This function is designed to make the framework compatible with the old plug-in, and the newly developed plug-in should no longer use this method
func (bot *Bot) OnGroupInvited(f func(qqClient *client.QQClient, event *client.GroupInvitedRequest)) {
	bot.GroupInvitedEvent.Subscribe(f)
}

// Deprecated: This function is designed to make the framework compatible with the old plug-in, and the newly developed plug-in should no longer use this method
func (bot *Bot) OnUserWantJoinGroup(f func(qqClient *client.QQClient, event *client.UserJoinGroupRequest)) {
	bot.UserWantJoinGroupEvent.Subscribe(f)
}

// Deprecated: This function is designed to make the framework compatible with the old plug-in, and the newly developed plug-in should no longer use this method
func (bot *Bot) OnNewFriend(f func(qqClient *client.QQClient, event *client.NewFriendEvent)) {
	bot.NewFriendEvent.Subscribe(f)
}

// Deprecated: This function is designed to make the framework compatible with the old plug-in, and the newly developed plug-in should no longer use this method
func (bot *Bot) OnNewFriendRequest(f func(qqClient *client.QQClient, event *client.NewFriendRequest)) {
	bot.NewFriendRequestEvent.Subscribe(f)
}

// Deprecated: This function is designed to make the framework compatible with the old plug-in, and the newly developed plug-in should no longer use this method
func (bot *Bot) OnDisconnected(f func(qqClient *client.QQClient, event *client.ClientDisconnectedEvent)) {
	bot.DisconnectedEvent.Subscribe(f)
}

// Deprecated: This function is designed to make the framework compatible with the old plug-in, and the newly developed plug-in should no longer use this method
func (bot *Bot) OnGroupNotify(f func(qqClient *client.QQClient, event client.INotifyEvent)) {
	bot.GroupNotifyEvent.Subscribe(f)
}

// Deprecated: This function is designed to make the framework compatible with the old plug-in, and the newly developed plug-in should no longer use this method
func (bot *Bot) OnFriendNotify(f func(qqClient *client.QQClient, event client.INotifyEvent)) {
	bot.FriendNotifyEvent.Subscribe(f)
}

// Deprecated: This function is designed to make the framework compatible with the old plug-in, and the newly developed plug-in should no longer use this method
func (bot *Bot) OnMemberSpecialTitleUpdated(f func(qqClient *client.QQClient, event *client.MemberSpecialTitleUpdatedEvent)) {
	bot.MemberSpecialTitleUpdatedEvent.Subscribe(f)
}

// Deprecated: This function is designed to make the framework compatible with the old plug-in, and the newly developed plug-in should no longer use this method
func (bot *Bot) OnGroupDigest(f func(qqClient *client.QQClient, event *client.GroupDigestEvent)) {
	bot.GroupDigestEvent.Subscribe(f)
}

// Deprecated: This function is designed to make the framework compatible with the old plug-in, and the newly developed plug-in should no longer use this method
func (bot *Bot) OnOtherClientStatusChanged(f func(qqClient *client.QQClient, event *client.OtherClientStatusChangedEvent)) {
	bot.OtherClientStatusChangedEvent.Subscribe(f)
}

// Deprecated: This function is designed to make the framework compatible with the old plug-in, and the newly developed plug-in should no longer use this method
func (bot *Bot) OnOfflineFile(f func(qqClient *client.QQClient, event *client.OfflineFileEvent)) {
	bot.OfflineFileEvent.Subscribe(f)
}
