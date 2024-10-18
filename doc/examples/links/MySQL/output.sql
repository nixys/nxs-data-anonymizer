-- MySQL dump 10.13  Distrib 8.4.0, for Linux (x86_64)
--
-- Host: localhost    Database: name_db
-- ------------------------------------------------------
-- Server version	8.4.0

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!50503 SET NAMES utf8mb4 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

--
-- Table structure for table `authors`
--

DROP TABLE IF EXISTS `authors`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `authors` (
  `id` int NOT NULL AUTO_INCREMENT,
  `first_name` varchar(50) COLLATE utf8mb3_unicode_ci NOT NULL,
  `last_name` varchar(50) COLLATE utf8mb3_unicode_ci NOT NULL,
  `email` varchar(100) COLLATE utf8mb3_unicode_ci NOT NULL,
  `birthdate` date NOT NULL,
  `added` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `email` (`email`)
) ENGINE=InnoDB AUTO_INCREMENT=11 DEFAULT CHARSET=utf8mb3 COLLATE=utf8mb3_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `authors`
--

LOCK TABLES `authors` WRITE;
/*!40000 ALTER TABLE `authors` DISABLE KEYS */;
INSERT INTO `authors` VALUES (26,'lwjNB8ofmQjcmzIKx23U','pBfm9pHABWxqgVMNk1pY','brett35@example.com','1999-12-31','2000-01-01 12:00:00'),(48,'xtc5AqJUOuAGGBir1Gvn','zdi0Y2ntYLIdTMGBhsRd','aracely46@example.com','1999-12-31','2000-01-01 12:00:00'),(9,'719GkwVKuXEKjAAohgrk','jU68LcKRCwZ3yFAOsQws','carroll.harris@example.com','1999-12-31','2000-01-01 12:00:00'),(15,'Bcync5tDYRvDOad9ulHY','70FuTHdJ2diuXKn5FLUl','johns.janick@example.org','1999-12-31','2000-01-01 12:00:00'),(49,'JsRPcdhoD7eUSG2yqaln','LR33Yhxv4bs2LdcvhTvg','lakin.ramiro@example.net','1999-12-31','2000-01-01 12:00:00'),(10,'jckw1fZr214qyZk2N5qd','Dx1dlwV6M9bownRFj8jv','judson33@example.com','1999-12-31','2000-01-01 12:00:00'),(46,'LynAXGoDTGMsFWxjggUr','mW9QmnaWOWeG0uTR3y4v','kaci.koch@example.net','1999-12-31','2000-01-01 12:00:00'),(8,'ZhCmwfvwrOmSTKN9vXrw','MuUlp8h02AyIQapcTNid','jprosacco@example.net','1999-12-31','2000-01-01 12:00:00'),(18,'04b4LKxXcOioU5JkO6Iy','SBoyCTw1RCtx0rMp6qJ7','kutch.kylie@example.com','1999-12-31','2000-01-01 12:00:00'),(42,'21L1p5ixBhOvI3P3iF5F','xsRxqdLYqLIGpdpNYrNA','hane.terrill@example.org','1999-12-31','2000-01-01 12:00:00');
/*!40000 ALTER TABLE `authors` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `posts`
--

DROP TABLE IF EXISTS `posts`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `posts` (
  `id` int NOT NULL AUTO_INCREMENT,
  `author_id` int NOT NULL,
  `title` varchar(255) COLLATE utf8mb3_unicode_ci NOT NULL,
  `description` varchar(500) COLLATE utf8mb3_unicode_ci NOT NULL,
  `content` text COLLATE utf8mb3_unicode_ci NOT NULL,
  `date` date NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=31 DEFAULT CHARSET=utf8mb3 COLLATE=utf8mb3_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `posts`
--

LOCK TABLES `posts` WRITE;
/*!40000 ALTER TABLE `posts` DISABLE KEYS */;
INSERT INTO `posts` VALUES (89,26,'Pariatur est in ut provident vero eligendi.','Consequatur quae odio in animi. Quis dolorem aut corporis sed ratione animi et. Iure sunt qui eligendi laudantium harum maxime id id. Quia eos beatae temporibus eos rerum ipsum.','Velit voluptatibus quo est dolorem quos sed et. Illum similique atque et qui. Perspiciatis magnam saepe quis. Consequatur aut dolorem ea.','1975-07-07'),(5,48,'Unde est architecto est.','Consectetur quidem aut eligendi quas. Voluptas aut voluptas ratione reiciendis magnam sint. Ratione et omnis assumenda deserunt voluptatem perspiciatis. Non voluptates quia distinctio nam consectetur vel.','Sed sint at consectetur. Fugiat fugiat omnis debitis dolore. Incidunt vero incidunt enim reiciendis. Odio et illo facilis quo sit quia consequatur voluptatem.','1976-05-25'),(56,9,'Dignissimos perspiciatis nobis quisquam saepe ad aut numquam.','Voluptatibus eum voluptas eveniet veniam ducimus. Modi id vitae est accusamus. Quaerat perferendis voluptas minima dignissimos et. Sapiente explicabo voluptas commodi voluptatem autem totam. Adipisci quos quos quo ea aliquid rerum.','Repellat minus enim occaecati sunt quas. Ab ipsa voluptatibus sunt eos omnis quisquam at voluptatibus. Quo vitae doloremque nesciunt doloremque. Temporibus rerum sunt iste dolores rerum.','1991-12-21'),(46,15,'Voluptas sint modi magnam.','Facere harum pariatur quo eveniet dolor molestias. Quo debitis corporis quasi minima et optio. Ut maiores nihil rerum autem culpa et voluptates. Amet totam rem optio in. Provident ea repudiandae quisquam unde occaecati autem.','Sed quia assumenda quis rerum praesentium harum. Quia minima quo natus. Impedit temporibus perspiciatis et enim doloremque nihil et. Perspiciatis nam occaecati illo dicta asperiores eos nam.','2002-01-23'),(76,49,'Sapiente rem eos enim ullam ipsum ut.','Officia quos voluptas autem consequuntur. Quis consequatur aut vel dolor. Sapiente voluptatem quibusdam dolor nobis earum laboriosam quisquam.','Quibusdam vel tempora explicabo qui. Sunt animi a unde veritatis perferendis similique nostrum vel. Porro ab ad cumque dolorum provident distinctio est ipsum. Esse maxime ea quia.','2017-11-19'),(99,10,'Error quas doloremque est sunt quae.','Corporis quis in sint. Eos vitae aut provident distinctio ullam. Corporis mollitia quam natus qui sapiente officiis.','Ipsam animi aperiam ipsa fuga. Sunt blanditiis possimus aut quod. Doloremque possimus occaecati numquam omnis dolor. Qui nihil qui atque vitae sapiente illo id.','1974-03-23'),(13,46,'Quo impedit quos molestiae dolorum in soluta dolores non.','Harum quasi dolorum et harum iste eveniet et. Cumque ab recusandae architecto est ipsam est. Eligendi et earum ea alias odio sed ut.','Illo earum porro corporis sunt aut. Vel quisquam ut voluptas reiciendis maiores cumque. Qui hic maiores voluptas sapiente reprehenderit. Eos non autem quis maiores perferendis.','2004-12-08'),(42,8,'Nulla quia repellendus et autem vitae provident.','Facere deserunt suscipit quia et et totam vel at. Cumque suscipit est ea quis in quos. Harum dolorum consequuntur illum voluptatem iste enim recusandae.','Excepturi hic tenetur sint nostrum. Nihil itaque repudiandae qui hic minima quis ut. Dolor et aut quaerat exercitationem minima ut. Blanditiis minus et qui dolor ut atque.','1999-02-28'),(61,18,'Similique at quia quia ut recusandae repudiandae delectus.','Et necessitatibus aut accusantium quisquam harum est. Non quam id impedit deleniti. Eaque non architecto facere dolorem nesciunt perferendis. Delectus nihil ut dignissimos. Voluptate rerum reprehenderit aut voluptatem quibusdam quidem.','Tempore iusto minus omnis tempore. Hic aut sequi temporibus consequatur. In vel enim eos nihil sed eos debitis.','2007-08-07'),(11,42,'Est vel aperiam ipsa quod doloremque et est et.','Facere ut et similique voluptates voluptas blanditiis explicabo. Eaque molestiae nihil sed et repellat voluptatem eum autem.','Rem laudantium in aut assumenda. Aspernatur id illo pariatur aut deleniti rem et. Nisi velit neque qui quia.','2009-07-18'),(87,26,'Sint est qui dolorem eum accusantium repudiandae.','Quia tenetur culpa maiores molestias. Id numquam illum earum quos sint ad dolore corrupti. Consequatur quasi itaque est odit qui quod culpa.','Omnis tenetur occaecati accusamus quis corrupti et ipsam. Ullam nobis tempore officia nesciunt iste. Nesciunt vel in eos. Dolor voluptates quod sed quibusdam ut. Fuga quia eum quidem cum.','2010-04-16'),(4,48,'Culpa debitis ut non sapiente voluptatem.','Commodi fugit ullam quaerat quam quo minus. Harum quam ipsam ducimus sit expedita sit. Eos natus quo quibusdam quam repudiandae. Assumenda et sint sit quia qui necessitatibus.','Eum molestiae cupiditate ut minus. Eaque eos eos ipsam voluptatem. Sint aut aliquid modi id dolores consequatur. Aut nemo blanditiis nisi ea nam velit.','1988-12-26'),(95,9,'Ut ut rerum qui quis sed sunt.','Tenetur enim et nisi ex et consequatur omnis nisi. Eos deleniti aut amet illo sed quibusdam. Consequuntur sint quia quia aut magni.','Aut eveniet natus ut. Nostrum aliquid qui nam hic iusto id occaecati qui. Est dolor accusamus qui saepe voluptatem dolorem omnis.','2001-06-01'),(52,15,'Error inventore delectus sapiente non.','Aliquam ut aliquid et laboriosam distinctio ab. Tempore molestiae aut ea nesciunt culpa. Eum maiores molestiae voluptatem animi aliquam. Aspernatur aut soluta totam occaecati in nam ratione asperiores.','Voluptatum id amet qui quia. Nemo nulla atque dignissimos.','1974-12-07'),(49,49,'Ut aut quo aut ea occaecati est voluptas.','Non deleniti molestias velit ut eos veritatis. Et id eaque numquam ipsa iure ut. Molestiae et minus natus.','Necessitatibus asperiores adipisci eaque. Alias sunt illum fugit ab labore laborum est odio. Aspernatur esse qui ut omnis autem. Laudantium quos corporis numquam sapiente enim.','1999-09-12'),(82,10,'Nihil ea qui sequi odit rem.','Omnis ea accusantium qui commodi et ducimus. Magni deserunt sed quo velit. Hic cum tenetur facilis. Non eum quis molestias sunt. Illo saepe architecto doloremque illum sequi.','Nihil iste illum vitae ipsa et et. Reiciendis sint illum quo dignissimos veritatis. Velit eum architecto enim ratione aut.','2003-07-28'),(40,46,'Quia accusantium deserunt suscipit.','Repellat et molestiae magnam quaerat et iure. Enim enim cum quisquam quasi.','Maiores officiis beatae hic hic dolor non officiis. Nulla beatae at aut fugit doloribus et provident. Debitis non dolores iusto impedit. Alias blanditiis velit voluptate et debitis.','2015-05-24'),(28,8,'Inventore et eaque temporibus qui aut exercitationem necessitatibus.','Omnis molestias nemo culpa est dicta dolorem porro. Illum est ut ut molestias similique sit sunt. Consequatur asperiores quaerat voluptate incidunt.','Unde quibusdam asperiores voluptates quaerat et aliquam. Libero at est quibusdam eos aut. Aut voluptatem tempora a nam.','1996-02-06'),(90,18,'Et soluta consequatur porro rem corrupti.','Ullam doloremque explicabo ea ab quam ipsum. Et qui iure vero tempore voluptate aut. Ullam id ab accusamus optio voluptas. Quia quaerat quibusdam aut accusantium impedit et.','Perspiciatis sed dolores aut. Animi sed totam exercitationem fugiat dolores. Ipsa dolorem quia non aliquam. Rerum exercitationem voluptatem dolor beatae veritatis quasi omnis consectetur.','2018-03-31'),(36,42,'Culpa error rerum voluptatem recusandae quae tempora.','Est ullam tempore autem et. Dolorum est eius eaque veniam in quibusdam neque. Ut aspernatur earum sint delectus voluptate deserunt sapiente quia.','Enim quos facere rerum ipsa. Eveniet expedita id eligendi voluptatem magnam quo eos.','2007-01-22'),(9,26,'Harum laudantium fugiat debitis atque sed.','Et asperiores magni qui voluptas a consequatur totam. Inventore enim quaerat consequatur sint quam voluptatibus optio. Est aut repellendus repudiandae reiciendis eveniet veniam. Minima a corrupti non quae.','Voluptate officia ut debitis explicabo corrupti facere reprehenderit. Enim id recusandae vitae est dolores. Aperiam consequuntur ex beatae et alias quia.','1987-06-29'),(77,48,'Non architecto ut voluptas aut voluptas dolor ut.','Dolore delectus optio fuga aut labore necessitatibus. Laborum quis nulla deserunt. Aut corporis repellendus inventore vel.','Facilis soluta ab culpa est. Aut pariatur et dicta sit. Sed qui est numquam est.\nAtque ex sit eveniet. Iure accusamus optio placeat voluptate non sit ea. Sint dignissimos dolorum debitis ipsam neque.','1993-05-06'),(18,9,'Dolores sed hic vitae ut qui.','Libero illum dolorum est eaque ut. Nesciunt qui vitae recusandae eveniet saepe et ut. Aliquid enim vitae beatae animi quam doloribus accusamus. Eos eos excepturi accusamus ut recusandae.','Distinctio molestiae sint qui. Expedita qui ex dignissimos architecto quo sapiente alias eveniet. Laborum qui aut est doloribus a beatae. Unde adipisci consequuntur et tenetur.','1978-09-08'),(29,15,'Deleniti distinctio eos eveniet.','Ex aperiam animi quis excepturi. Velit qui molestias iste iure voluptatem. Eveniet adipisci quasi amet. Soluta laborum aperiam magni.','Deserunt dolores ut ut ex cum. Quae iusto rerum voluptatum aperiam. Vel voluptatem ut nesciunt delectus blanditiis. Sequi consequuntur temporibus sunt dolores repudiandae.','2003-06-23'),(80,49,'Ea et asperiores odio vel sunt exercitationem.','Quod necessitatibus quis sit aspernatur a. Molestias laudantium voluptatibus quaerat rerum cupiditate. Aut consectetur quia doloribus repellat porro est.','Et alias soluta aliquam qui aut. Explicabo quo consequuntur consequuntur velit distinctio.','1972-09-16'),(69,10,'Accusamus est sed aut minus.','Et libero in in est consequuntur. Illum distinctio doloremque quas.','Neque vel tenetur libero sed aut voluptas. Itaque quod et laudantium magnam tempora qui non. Et sapiente dolores reiciendis sit ab fugiat sed sequi. Reiciendis maxime qui aut laborum dolor.','2017-03-08'),(98,46,'Et sed reprehenderit nobis sed.','Totam amet nisi id quod. Laudantium maxime rem velit consequuntur amet. Et enim eum aspernatur totam. Nulla consequatur similique maxime et qui rerum.','Iure sequi unde consectetur neque in magni. Illo ipsa commodi quae placeat modi numquam numquam. Quae culpa sit distinctio vitae est.','1973-08-04'),(27,8,'Autem ullam aperiam totam assumenda quod esse.','Aut incidunt aspernatur aut sapiente ipsum repellat. Molestiae eaque accusantium maxime. Velit modi ea consectetur natus ea.','Molestias voluptatum cupiditate et minus qui consequuntur. Reprehenderit aut at corrupti. Eius unde ullam aut cum eius molestiae.','2014-04-26'),(64,18,'Ad dolores amet ea nisi aut enim voluptatem.','Eum voluptates sunt occaecati molestiae. Cupiditate labore expedita eius in omnis non. Sit cumque unde cupiditate sequi eius reprehenderit nulla.','Voluptates modi esse quas sit eos praesentium. Quas et rem vero quo veritatis voluptas. Inventore dolorem ut praesentium consequatur rem. Amet ipsam quia officia iste est in.','2006-09-07'),(32,42,'Eos fuga at ex sapiente quasi.','Facere dolores non omnis facilis. Nam cupiditate maiores iure quia adipisci numquam magnam. Esse voluptas voluptatem cum ea voluptas doloremque.','Voluptatem qui qui magnam quis eos. Consequatur ut repudiandae libero dignissimos et quia velit. Eum nisi facere consequatur quod ut. Ad molestiae voluptatem eaque iste et.','2009-07-23');
/*!40000 ALTER TABLE `posts` ENABLE KEYS */;
UNLOCK TABLES;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2024-07-19 10:47:53
