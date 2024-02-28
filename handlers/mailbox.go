package handlers

import (
	"distriai-index-solana/common"
	"distriai-index-solana/model"
	"distriai-index-solana/utils/logs"
	"distriai-index-solana/utils/resp"
	"fmt"
	"github.com/gin-gonic/gin"
	"gopkg.in/gomail.v2"
)

const EMAIL_BODY = `
<div>
	<div style="line-height:1.7;color:#000000;font-size:14px;font-family:Arial">
		<div>
			<div style="font-family: Helvetica, &quot;Microsoft Yahei&quot;, verdana; position: relative; word-break: break-word !important;">
				<div style="position: relative; background-color: rgb(247, 247, 247); word-break: break-word !important;">
					<img src="https://distriai.cloud/logo2.png" style="border: 0px; width: 200px; margin: 0px auto; display: block; padding-top: 100px; word-break: break-word !important;" />
					<table width="100%%" cellpadding="0" cellspacing="0" border="0" style="font-size: 14px; font-family: &quot;Microsoft Yahei&quot;, Arial, Helvetica, sans-serif; padding: 0px; margin: 0px; color: rgb(51, 51, 51); background-repeat: repeat-x; background-position: left bottom; word-break: break-word !important; border-collapse: collapse; border-spacing: 0px;" id="ntes_editor_table_10010">
						<tbody style="word-break: break-word !important;">
							<tr style="word-break: break-word !important;">
								<td style="font-family: arial, verdana, sans-serif; line-height: 1.666; word-break: break-word !important;">
									<table width="600" border="0" align="center" cellpadding="0" cellspacing="0" style="word-break: break-word !important; border-collapse: collapse; border-spacing: 0px;" id="ntes_editor_table_10011">
										<tbody style="word-break: break-word !important;">
											<tr style="word-break: break-word !important;">
												<td align="center" valign="middle" style="line-height: 1.666; padding: 0px; word-break: break-word !important;">
													<br style="word-break: break-word !important;" />
													<br style="word-break: break-word !important;" />
													<br style="word-break: break-word !important;" />
													<br style="word-break: break-word !important;" />
												</td>
											</tr>
											<tr style="word-break: break-word !important;">
												<td style="line-height: 1.666; word-break: break-word !important;">
													<div style="padding: 0px 30px; background: rgb(255, 255, 255); box-shadow: rgb(238, 238, 238) 0px 0px 5px; word-break: break-word !important;">
														<table width="100%%" border="0" cellspacing="0" cellpadding="0" style="word-break: break-word !important; border-collapse: collapse; border-spacing: 0px;" id="ntes_editor_table_10012">
															<tbody style="word-break: break-word !important;">
																<tr style="word-break: break-word !important;">
																	<td style="line-height: 1.666; border-bottom: 1px solid rgb(230, 230, 230); font-size: 18px; padding: 20px 0px; word-break: break-word !important;">
																		<br style="word-break: break-word !important;" />
																		<table border="0" cellspacing="0" cellpadding="0" width="100%%" style="word-break: break-word !important; border-collapse: collapse; border-spacing: 0px;" id="ntes_editor_table_10013">
																			<tbody style="word-break: break-word !important;">
																				<tr style="word-break: break-word !important;">
																					<td style="line-height: 1.666; word-break: break-word !important;">
																						<font size="5" style="word-break: break-word !important;">Welcome to DistriAI ！
																							<br style="word-break: break-word !important;" />
																							<br style="word-break: break-word !important;" />
																						</font>
																					</td>
																					<td style="line-height: 1.666; word-break: break-word !important;"></td>
																				</tr>
																			</tbody>
																		</table>
																	</td>
																</tr>
																<tr style="word-break: break-word !important;">
																	<td style="line-height: 30px; padding: 20px 0px 0px; word-break: break-word !important;">
																		<p dir="ltr" style="font-size: 14px; color: rgb(102, 102, 102); word-break: break-word !important;">Thanks for Subscribe DistriAI. DistriAI aims to establish a fair, efficient, and transparent AI computing power network to meet the growing demand for computing power and reduce its usage costs. Currently, DistriAI is in the testing phase and is expected to officially launch in November.</p>
																		<p dir="ltr" style="word-break: break-word !important;">
																			<spanmicrosoft yahei="" style="color: rgb(102, 102, 102); font-size: 14px; font-style: italic; word-break: break-word !important;">If you have any questions, please contact us via&nbsp;<a href="https://twitter.com/DistriAI_web3" style="color: rgb(51, 204, 204); word-break: break-word !important;">Twitter&nbsp;</a>&nbsp;or&nbsp;<a href="https://t.me/hanleeeeeee" style="color: rgb(51, 204, 204); word-break: break-word !important;">Telegram&nbsp;</a>. We welcome your suggestions and feedback.</spanmicrosoft>
																		</p>
																		<p dir="ltr" style="word-break: break-word !important;">
																			<span style="color: rgb(102, 102, 102); font-size: 14px; font-style: italic; word-break: break-word !important;">Friendly reminder: Please visit the correct DistriAI address at&nbsp;</span>
																			<span style="font-size: 14px; word-break: break-word !important;">
																				<span style="font-style: italic; color: rgb(51, 204, 204); word-break: break-word !important;">
																					<a href="https://www.distriai.cloud/" style="color: rgb(51, 112, 255); word-break: break-word !important;">www.distriai.cloud</a>
																				</span>
																			</span>
																		</p>
																		<p dir="ltr" style="font-size: 14px; color: rgb(102, 102, 102); word-break: break-word !important;">
																			<spanmicrosoft yahei="" style="font-style: italic; word-break: break-word !important;">Thanks for your support</spanmicrosoft>
																		</p>
																	</td>
																</tr>
																<tr style="word-break: break-word !important;">
																	<td style="line-height: 20px; padding: 30px 0px 15px; font-size: 12px; color: rgb(153, 153, 153); word-break: break-word !important;">DistriAI Team
																		<br style="word-break: break-word !important;" />System email, please do not reply.</td>
																</tr>
															</tbody>
														</table>
													</div>
												</td>
											</tr>
											<tr style="word-break: break-word !important;">
												<td align="center" style="line-height: 1.666; padding: 26px 0px; word-break: break-word !important;">
													<br style="word-break: break-word !important;" />
												</td>
											</tr>
											<tr style="word-break: break-word !important;">
												<td align="center" style="line-height: 1.666; word-break: break-word !important;">
													<br style="word-break: break-word !important;" />
													<span style="color: rgb(153, 153, 153); font-size: 12px; word-break: break-word !important;">If you do not wish to receive this type of email,&nbsp;</span>
													<a href="https://www.distriai.cloud/mailbox/unsubscribe/%s" style="font-size: 12px; color: rgb(51, 204, 204); word-break: break-word !important;">
														unsubscribe
													</a>
													<br style="word-break: break-word !important;" />
												</td>
											</tr>
											<tr style="word-break: break-word !important;">
												<td align="center" style="line-height: 1.666; font-size: 12px; color: rgb(153, 153, 153); padding: 5px 0px 20px; word-break: break-word !important;">© 2023 DistriAI All Rights Reserved
													<br style="word-break: break-word !important;" />Official website：&nbsp;www.distri.ai
													<br style="word-break: break-word !important;" />
													<br style="word-break: break-word !important;" />
													<br style="word-break: break-word !important;" />
													<br style="word-break: break-word !important;" />
													<br style="word-break: break-word !important;" />
												</td>
											</tr>
										</tbody>
									</table>
								</td>
							</tr>
						</tbody>
					</table>
				</div>
			</div>
			<br />
		</div>
	</div>
	<br />
</div>
`

type MailboxReq struct {
	Mailbox string `json:"mailbox" binding:"required"`
}

func Subscribe(context *gin.Context) {
	var req MailboxReq
	err := context.ShouldBindJSON(&req)
	if err != nil {
		logs.Warn(fmt.Sprintf("Subscribe paramter missing: %s", err))
		resp.Fail(context, "Parameter missing")
		return
	}

	mailbox := &model.Mailbox{MailBox: req.Mailbox}
	var count int64
	dbResult := common.Db.Model(&mailbox).Where(&mailbox).Count(&count)
	if dbResult.Error != nil {
		logs.Error(fmt.Sprintf("Subscribe database error: %s", dbResult.Error))
		resp.Fail(context, "Database error")
		return
	}
	if count > 0 {
		resp.Fail(context, "Mailbox already subscribed")
		return
	}

	dbResult = common.Db.Create(&mailbox)
	if dbResult.Error != nil {
		logs.Error(fmt.Sprintf("Subscribe database error: %s", dbResult.Error))
		resp.Fail(context, "Database error")
		return
	}

	err = sendEmail(req.Mailbox, fmt.Sprintf(EMAIL_BODY, mailbox.MailBox))
	if err != nil {
		logs.Error(fmt.Sprintf("Send email Fail,error: %s", err))
		resp.Fail(context, "Send email Fail")
		return
	}

	resp.Success(context, "")
}

func Unsubscribe(context *gin.Context) {
	var req MailboxReq
	err := context.ShouldBindJSON(&req)
	if err != nil {
		logs.Warn(fmt.Sprintf("Unsubscribe paramter missing: %s", err))
		resp.Fail(context, "Parameter missing")
		return
	}

	mailbox := &model.Mailbox{MailBox: req.Mailbox}
	var count int64
	dbResult := common.Db.Model(&mailbox).Where(&mailbox).Count(&count)
	if dbResult.Error != nil {
		logs.Error(fmt.Sprintf("Unsubscribe database error: %s", dbResult.Error))
		resp.Fail(context, "Database error")
		return
	}
	if count == 0 {
		resp.Fail(context, "Mailbox is not subscribed")
		return
	}

	dbResult = common.Db.Where(&mailbox).Delete(&mailbox)
	if dbResult.Error != nil {
		logs.Error(fmt.Sprintf("Unsubscribe database error: %s", dbResult.Error))
		resp.Fail(context, "Database error")
		return
	}

	resp.Success(context, "")
}

func sendEmail(to string, text string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", common.Conf.Mailbox.Username)
	m.SetHeader("To", to)
	m.SetHeader("Subject", "Welcome to DistriAI")
	m.SetBody("text/html", text)

	d := gomail.NewDialer(
		common.Conf.Mailbox.Host,
		common.Conf.Mailbox.Port,
		common.Conf.Mailbox.Username,
		common.Conf.Mailbox.Password,
	)
	err := d.DialAndSend(m)
	if err != nil {
		fmt.Println(err)
	}
	return err
}
