package dao

import (
	"database/sql"
	"errors"
	"fmt"
	"time"
	"winter-test/model"
)

// 用户新建立评论
func CreateComment(username string, productID int, comment string, parentID int) (*model.Comment, error) {
	// 判断用户是否存在
	var commentAll model.Comment

	err := db.QueryRow("select nickname from user where username = ?", username).Scan(&commentAll.Nickname)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("该用户不存在")
		}
		return nil, err
	}

	// 判断评论的商品是否存在
	err = db.QueryRow("select name from product where id = ?", productID).Scan(&commentAll.ProductName)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("该商品不存在")
		}
		return nil, err
	}

	// 将评论信息插入数据库
	res, err := db.Exec("insert into comments (user_name, product_id, comment, time, parent_id ) values (?, ?, ?, ?, ?)", username, productID, comment, time.Now(), parentID)
	if err != nil {
		return nil, err
	}

	// 获取评论的自增ID
	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}
	commentAll.Username = username
	db.QueryRow("select * from comments where id = ?", id).Scan(&commentAll.ID, &commentAll.Username, &commentAll.ProductID, &commentAll.Comment, &commentAll.Time, &commentAll.ParentID)
	return &commentAll, nil
}

// 展示商品的评论区
func GetComment(productID, parentID int) ([]model.NewComment, error) {
	var comments []model.NewComment
	fmt.Println(productID, parentID)
	rows, err := db.Query(`select comments.id, comments.user_name, user.NickName, comments.product_id, product.name, comments.comment, comments.time, comments.parent_id
from comments
join user ON comments.user_name = user.UserName
join product ON comments.product_id = product.ID
where comments.product_id = ? AND comments.parent_id = ?;`, productID, parentID)
	if err != nil {
		return nil, err

	}

	defer rows.Close()
	for rows.Next() {
		var comment model.Comment
		err = rows.Scan(&comment.ID, &comment.Username, &comment.Nickname, &comment.ProductID, &comment.ProductName, &comment.Comment, &comment.Time, &comment.ParentID)
		fmt.Println(comment, err, 3333)
		if err != nil {

			return nil, err
		}
		fmt.Println(comment)
		// 查出子评论数量
		rows2, err1 := db.Query("select * from comments where parent_id = ? and product_id = ?", comment.ID, productID)
		if err1 != nil {
			return nil, err1
		}
		var count int
		for rows2.Next() {
			count++
		}

		rows2.Scan(&count)
		newComment := model.NewComment{
			Comment:  comment,
			SonCount: count,
		}

		fmt.Println(comment, 11111111111)
		comments = append(comments, newComment)
		fmt.Println(comments, 00000000000)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return comments, nil
}

// 查询评论是否属于该用户
func GetCommentByID(commentID int) (string, error) {
	var username string
	if err := db.QueryRow("select user_name from comments where id = ?", commentID).Scan(&username); err != nil {
		return username, err
	}
	return username, nil
}

// 删除评论
func DeleteComment(commentID int) error {
	_, err := db.Exec("delete from comments where id = ?", commentID)
	if err != nil {
		return err
	}

	return nil
}
