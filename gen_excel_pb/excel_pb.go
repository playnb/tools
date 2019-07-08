package main

import (
	"fmt"
	//"common/mustang/tbl"
	"cell/conf"
	rdExcel "cell/conf/tbl"
	"cell/common/protocol/tbx"
	"flag"
	"os"
)

var ConfigDir string

func init() {
}

func main() {
	flag.Set("config", "./")
	//flag.Set("config", `C:\code\doc\`)
	flag.Parse()
	fmt.Println(conf.GetMe().ConfigDir)
	rdExcel.ReadTbl()
	rdExcel.LoadConfig()
	fmt.Println("读取Excel文件")

	tbxFile := &tbx.TbxFile{}

	//技能表
	tbxFile.Skills = make(map[uint64]*tbx.TbxSkillBase)
	rdExcel.ForeachSkillTbl(func(skill *rdExcel.SkillTbl) {
		tbxFile.Skills[skill.SkillID] = &tbx.TbxSkillBase{
			SkillID:       skill.SkillID,
			GravityEffect: skill.GravityEffect,
			WindEffect:    skill.WindEffect,
			Mass:          skill.Mass,
			SpeedEffect:   skill.SpeedEfect,
			BulletRadius:  skill.BulletRadius,
		}
	})

	//英雄表
	tbxFile.Heros = make(map[uint64]*tbx.TbxHeroBase)
	rdExcel.ForeachHeroTbl(func(hero *rdExcel.HeroTbl) {
		tbxFile.Heros[hero.HeroID] = &tbx.TbxHeroBase{
			HeroID: hero.HeroID,
			Skill1: hero.Skill1[0],
			Skill2: hero.Skill2[0],
			Skill3: hero.Skill3[0],
			Skill4: hero.Skill4[0],
		}
	})

	fmt.Println("==============================> 写入Proto文件...")
	data, err := tbxFile.Marshal()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("data len: %d\n", len(data))
	f, err := os.OpenFile(conf.GetMe().ConfigDir+"/Excel/tbx.pb", os.O_CREATE|os.O_TRUNC, os.ModePerm)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()
	f.Write(data)
}
