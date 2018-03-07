package smsTool

import (
    "testing"
    "work.goproject.com/goutil/smsUtil/qcloud"
)

func TestQCloudSend(t *testing.T) {
    sms, err := qcloud.NewTmplSms("123", "abcdef",
                        []string{"86"},
                        []string{"13100000000"},
                        14176, []string{"9527"},
                        "摩奇互娱", "", "hello-world!")
    if err != nil {
        t.Errorf("qcloud.TmplSms error: %v", err)
        return
    }

    if ok, err := sms.Send(); !ok {
        if err != nil {
            t.Errorf("qcloud.Send error: %v", err)
            return
        } else if !ok {
            rspn := sms.GetResponse().(*qcloud.QCloudResponse)
            t.Errorf("qcloud return error: %v", rspn)
        }
    }

    sms, err = qcloud.NewMsgSms("123", "abcdefg",
                        []string{"86"},
                        []string{"13100000000"},
                        "【摩奇互娱】尊敬的用户您的短信验证码为123456，验证码时间为10分钟内有效，请尽快完成填报，感谢您的支持",
                        "", "hello-world!")
    if err != nil {
        t.Errorf("qcloud.TmplSms error: %v", err)
        return
    }

    if ok, err := sms.Send(); !ok {
        if err != nil {
            t.Errorf("qcloud.Send error: %v", err)
            return
        } else if !ok {
            rspn := sms.GetResponse().(*qcloud.QCloudResponse)
            t.Errorf("qcloud return error: %v", rspn)
        }
    }
}

