//	Copyright 2025 Naveen R
//
//		Licensed under the Apache License, Version 2.0 (the "License");
//		you may not use this file except in compliance with the License.
//		You may obtain a copy of the License at
//
//		http://www.apache.org/licenses/LICENSE-2.0
//
//		Unless required by applicable law or agreed to in writing, software
//		distributed under the License is distributed on an "AS IS" BASIS,
//		WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//		See the License for the specific language governing permissions and
//		limitations under the License.
package com.naveen_r_sam.auth_service.common;

import org.springframework.amqp.rabbit.core.RabbitTemplate;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;

@Service
public class MessageSender {
    private final RabbitTemplate rabbitTemplate;
    private final String defaultQueue;

    // Constructor with default queue name
    @Autowired
    public MessageSender(RabbitTemplate rabbitTemplate) {
        this.rabbitTemplate = rabbitTemplate;
        this.defaultQueue = "error_dump_queue"; // Set your default queue name here
    }

    public MessageSender(RabbitTemplate rabbitTemplate, String queueName) {
        this.rabbitTemplate = rabbitTemplate;
        this.defaultQueue = queueName; // Set the queue name provided in the constructor
    }

    // Send message to a specific queue
    public void sendMessage(String queueName, Object message) {
        rabbitTemplate.convertAndSend(queueName, message);
        System.out.println("Message sent to queue [" + queueName + "]: " + message);
    }

    // Send message to the default queue
    public void sendMessage(Object message) {
        rabbitTemplate.convertAndSend(defaultQueue, message);
        System.out.println("Message sent to default queue [" + defaultQueue + "]: " + message);
    }
}