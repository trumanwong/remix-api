CREATE TABLE `tasks` (
                         `id` bigint unsigned NOT NULL AUTO_INCREMENT,
                         `user_id` bigint unsigned NOT NULL COMMENT '用户id',
                         `uuid` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'uuid',
                         `serial` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '任务号',
                         `status` tinyint unsigned NOT NULL COMMENT '状态0处理中，1处理成功，2处理失败',
                         `file_name` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '文件名',
                         `size` bigint unsigned NOT NULL COMMENT '文件大小',
                         `path` varchar(160) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '保存路径',
                         `convert_path` varchar(160) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '转换后存储路径',
                         `ip` bigint NOT NULL COMMENT '访问ip',
                         `deleted_at` timestamp NULL DEFAULT NULL,
                         `created_at` timestamp NULL DEFAULT NULL,
                         `updated_at` timestamp NULL DEFAULT NULL,
                         `type` tinyint unsigned NOT NULL DEFAULT '0' COMMENT '任务类型 0gif倒放 1诺基亚短信',
                         `other` json NOT NULL COMMENT '其他字段',
                         PRIMARY KEY (`id`),
                         UNIQUE KEY `uni_user_serial` (`user_id`,`serial`),
                         UNIQUE KEY `uni_uuid_serial` (`uuid`,`serial`)
) ENGINE=InnoDB AUTO_INCREMENT=110 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE `configs` (
                           `id` bigint unsigned NOT NULL AUTO_INCREMENT,
                           `type` tinyint unsigned NOT NULL COMMENT '0鬼畜动图',
                           `content` json NOT NULL COMMENT '配置内容',
                           `created_at` timestamp NULL DEFAULT NULL,
                           `updated_at` timestamp NULL DEFAULT NULL,
                           PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;