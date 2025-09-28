-- 提示词相关数据表DDL

-- 分类表
CREATE TABLE `categories` (
    `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '分类ID',
    `name` varchar(100) NOT NULL COMMENT '分类名称',
    `description` text COMMENT '分类描述',
    `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_name` (`name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='提示词分类表';

-- 提示词表
CREATE TABLE `prompts` (
    `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '提示词ID',
    `title` varchar(200) NOT NULL COMMENT '标题',
    `category_id` int(11) NOT NULL COMMENT '分类ID',
    `prompt` text NOT NULL COMMENT '提示词内容',
    `image_description` text COMMENT '图片描述',
    `oss_short_links` json NOT NULL COMMENT 'OSS短链列表，JSON数组格式',
    `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    KEY `idx_category_id` (`category_id`),
    KEY `idx_title` (`title`),
    CONSTRAINT `fk_prompts_category` FOREIGN KEY (`category_id`) REFERENCES `categories` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='提示词表';