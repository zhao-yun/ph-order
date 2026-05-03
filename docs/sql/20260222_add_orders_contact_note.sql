ALTER TABLE `orders`
  ADD COLUMN `contact` VARCHAR(255) NULL AFTER `sitter_name`,
  ADD COLUMN `alternative_contact` VARCHAR(255) NULL AFTER `contact`,
  ADD COLUMN `note` TEXT NULL AFTER `alternative_contact`;

