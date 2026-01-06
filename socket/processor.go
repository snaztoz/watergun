package socket

import "errors"

func newMessageProcessor() *messageProcessor {
	return &messageProcessor{
		clientToUserIDMapping: make(map[string]string),
	}
}

type messageProcessor struct {
	clientToUserIDMapping map[string]string
}

func (mp *messageProcessor) process(c *client, msg *ReadMessage) error {
	if mp.shouldIgnoreMessage(c, msg) {
		return nil
	}

	switch {
	case msg.isForAuthentication():
		return mp.authenticate(c, msg)
	}

	return nil
}

func (mp *messageProcessor) authenticate(c *client, msg *ReadMessage) error {
	if msg.Content != "SOME-CLIENT-SECRET-CHANGE-LATER" {
		return errors.New("invalid authentication credential")
	}

	mp.clientToUserIDMapping[c.id] = "some-user-id"

	return nil
}

func (mp *messageProcessor) shouldIgnoreMessage(c *client, msg *ReadMessage) bool {
	return mp.isRedundantAuthentication(c, msg) || mp.isUnauthenticatedMessage(c, msg)
}

func (mp *messageProcessor) isRedundantAuthentication(c *client, msg *ReadMessage) bool {
	return mp.isAuthenticated(c) && msg.isForAuthentication()
}

func (mp *messageProcessor) isUnauthenticatedMessage(c *client, msg *ReadMessage) bool {
	return !mp.isAuthenticated(c) && !msg.isForAuthentication()
}

func (mp *messageProcessor) isAuthenticated(c *client) bool {
	_, isAuthenticated := mp.clientToUserIDMapping[c.id]
	return isAuthenticated
}
